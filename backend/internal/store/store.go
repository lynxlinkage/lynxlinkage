// Package store wraps database access. It exposes a small repository per
// domain and lives behind interfaces, currently backed by PostgreSQL via
// jackc/pgx.
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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Open opens a PostgreSQL handle and runs embedded migrations.
//
// dsn is a libpq-style URL, e.g.
// postgresql://user:pass@host:5432/dbname?sslmode=disable.
func Open(ctx context.Context, dsn string, logger *slog.Logger) (*sqlx.DB, error) {
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	sqlDB := stdlib.OpenDB(*cfg)
	db := sqlx.NewDb(sqlDB, "pgx")
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	if err := runMigrations(ctx, db, logger); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

// runMigrations applies the SQL files embedded under migrations/. We use
// a minimal homegrown runner (no goose dependency at runtime) keyed off
// the filename; each file is executed in a single transaction with one
// statement per Exec so the extended pg protocol is happy.
func runMigrations(ctx context.Context, db *sqlx.DB, logger *slog.Logger) error {
	if _, err := db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version    TEXT PRIMARY KEY,
            applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
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
		stmts := splitStatements(up)

		if err := withTx(ctx, db, func(tx *sqlx.Tx) error {
			for _, stmt := range stmts {
				if _, err := tx.ExecContext(ctx, stmt); err != nil {
					return fmt.Errorf("exec %s: %w", name, err)
				}
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO schema_migrations(version) VALUES ($1)`, version); err != nil {
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

// splitStatements breaks a SQL blob into individual statements on
// top-level semicolons, ignoring delimiters that fall inside single- or
// double-quoted strings, line comments, or block comments. The
// PostgreSQL extended-query protocol used by pgx only accepts one
// statement per Exec, so the migration runner needs to feed each one
// separately.
func splitStatements(sqlText string) []string {
	var (
		out  []string
		buf  strings.Builder
		i    int
		n    = len(sqlText)
		cur  byte
		next byte
	)

	flush := func() {
		s := strings.TrimSpace(buf.String())
		if s != "" {
			out = append(out, s)
		}
		buf.Reset()
	}

	for i < n {
		cur = sqlText[i]
		if i+1 < n {
			next = sqlText[i+1]
		} else {
			next = 0
		}

		switch cur {
		case '\'', '"':
			quote := cur
			buf.WriteByte(cur)
			i++
			for i < n {
				c := sqlText[i]
				buf.WriteByte(c)
				if c == quote {
					i++
					// Handle '' / "" escapes inside string literals.
					if i < n && sqlText[i] == quote {
						buf.WriteByte(sqlText[i])
						i++
						continue
					}
					break
				}
				i++
			}
		case '-':
			if next == '-' {
				for i < n && sqlText[i] != '\n' {
					buf.WriteByte(sqlText[i])
					i++
				}
			} else {
				buf.WriteByte(cur)
				i++
			}
		case '/':
			if next == '*' {
				buf.WriteByte(cur)
				buf.WriteByte(next)
				i += 2
				for i < n {
					if sqlText[i] == '*' && i+1 < n && sqlText[i+1] == '/' {
						buf.WriteByte('*')
						buf.WriteByte('/')
						i += 2
						break
					}
					buf.WriteByte(sqlText[i])
					i++
				}
			} else {
				buf.WriteByte(cur)
				i++
			}
		case ';':
			flush()
			i++
		default:
			buf.WriteByte(cur)
			i++
		}
	}
	flush()
	return out
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
