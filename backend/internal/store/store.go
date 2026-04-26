// Package store wraps database access. It exposes a small repository per
// domain and lives behind interfaces so the migration to Postgres later is
// only a wiring change.
package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Open opens a database handle and runs embedded migrations.
//
// The driver name passed to sqlx must match the imported driver. We use the
// pure-Go modernc.org/sqlite driver (registered as "sqlite") which avoids
// CGO and produces a fully static binary.
func Open(ctx context.Context, dsn string, logger *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite serialises writes; one conn avoids "database is locked".
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := runMigrations(ctx, db, logger); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

// runMigrations applies the SQL files embedded under migrations/. We use a
// minimal homegrown runner (no goose dependency at runtime) keyed off the
// filename; each file is executed in a single transaction.
func runMigrations(ctx context.Context, db *sqlx.DB, logger *slog.Logger) error {
	if _, err := db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY,
            applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        )`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	applied := map[string]bool{}
	rows, err := db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}

	for _, name := range files {
		version := strings.TrimSuffix(name, ".sql")
		if applied[version] {
			continue
		}
		raw, err := migrationsFS.ReadFile(filepath.ToSlash(filepath.Join("migrations", name)))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		// Take only the "up" portion: everything after `-- +goose Up` and
		// before `-- +goose Down`. Strip statement begin/end markers.
		up := extractUp(string(raw))

		if err := withTx(ctx, db, func(tx *sqlx.Tx) error {
			if _, err := tx.ExecContext(ctx, up); err != nil {
				return fmt.Errorf("exec %s: %w", name, err)
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO schema_migrations(version) VALUES (?)`, version); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		logger.Info("migration applied", "version", version)
	}
	return nil
}

func extractUp(raw string) string {
	const upMarker = "-- +goose Up"
	const downMarker = "-- +goose Down"
	if i := strings.Index(raw, upMarker); i >= 0 {
		raw = raw[i+len(upMarker):]
	}
	if i := strings.Index(raw, downMarker); i >= 0 {
		raw = raw[:i]
	}
	raw = strings.ReplaceAll(raw, "-- +goose StatementBegin", "")
	raw = strings.ReplaceAll(raw, "-- +goose StatementEnd", "")
	return raw
}

func withTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// IsNoRows reports whether err is sql.ErrNoRows.
func IsNoRows(err error) bool { return err == sql.ErrNoRows }

// ErrNotFound is returned by repositories when an update target row does
// not exist.
var ErrNotFound = sql.ErrNoRows
