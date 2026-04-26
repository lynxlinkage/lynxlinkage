package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// StatusRepo persists HR-defined application statuses. The hiring
// pipeline graph itself is implicit: rows are sorted by display_order
// and HR moves applications between them freely. Two terminal kinds
// (`accept` / `reject`) tell the rest of the system which states are
// "closed" for filtering and metrics.
type StatusRepo struct{ db *sqlx.DB }

func NewStatusRepo(db *sqlx.DB) *StatusRepo { return &StatusRepo{db: db} }

const statusColumns = `id, slug, name, kind, color, display_order, is_default, created_at`

// ErrStatusInUse is returned by Delete when there are still
// applications pointing at the row.
var ErrStatusInUse = errors.New("store: status still has applications")

// ErrStatusNoDefault is returned by Update when an HR user tries to
// unset the last default status without nominating another.
var ErrStatusNoDefault = errors.New("store: at least one status must be default")

// List returns every status ordered for display.
func (r *StatusRepo) List(ctx context.Context) ([]domain.ApplicationStatus, error) {
	out := []domain.ApplicationStatus{}
	q := `SELECT ` + statusColumns + ` FROM application_statuses ORDER BY display_order ASC, id ASC`
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *StatusRepo) Get(ctx context.Context, id int64) (*domain.ApplicationStatus, error) {
	var s domain.ApplicationStatus
	q := `SELECT ` + statusColumns + ` FROM application_statuses WHERE id = ?`
	if err := r.db.GetContext(ctx, &s, q, id); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StatusRepo) GetBySlug(ctx context.Context, slug string) (*domain.ApplicationStatus, error) {
	var s domain.ApplicationStatus
	q := `SELECT ` + statusColumns + ` FROM application_statuses WHERE slug = ?`
	if err := r.db.GetContext(ctx, &s, q, slug); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetDefault returns the row currently marked is_default = 1. If for
// some reason no row carries the flag (HR may have edited carelessly),
// it falls back to the lowest display_order so the system can still
// label new applications.
func (r *StatusRepo) GetDefault(ctx context.Context) (*domain.ApplicationStatus, error) {
	var s domain.ApplicationStatus
	q := `
        SELECT ` + statusColumns + `
        FROM application_statuses
        WHERE is_default = 1
        LIMIT 1
    `
	err := r.db.GetContext(ctx, &s, q)
	if err == nil {
		return &s, nil
	}
	if !IsNoRows(err) {
		return nil, err
	}
	q2 := `
        SELECT ` + statusColumns + `
        FROM application_statuses
        ORDER BY display_order ASC, id ASC
        LIMIT 1
    `
	if err := r.db.GetContext(ctx, &s, q2); err != nil {
		return nil, err
	}
	return &s, nil
}

// Create inserts a new status. If is_default is true, any existing
// default is cleared in the same transaction.
func (r *StatusRepo) Create(ctx context.Context, s *domain.ApplicationStatus) (int64, error) {
	if s.Kind == "" {
		s.Kind = domain.StatusKindOpen
	}
	var id int64
	err := withTx(ctx, r.db, func(tx *sqlx.Tx) error {
		if s.IsDefault {
			if _, err := tx.ExecContext(ctx,
				`UPDATE application_statuses SET is_default = 0 WHERE is_default = 1`); err != nil {
				return err
			}
		}
		res, err := tx.ExecContext(ctx, `
            INSERT INTO application_statuses (slug, name, kind, color, display_order, is_default)
            VALUES (?, ?, ?, ?, ?, ?)`,
			s.Slug, s.Name, s.Kind, s.Color, s.DisplayOrder, s.IsDefault)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	if err != nil {
		return 0, err
	}
	s.ID = id
	return id, nil
}

// Update modifies a status. Promoting a row to default demotes the
// existing default; demoting the last default without promoting another
// returns ErrStatusNoDefault.
func (r *StatusRepo) Update(ctx context.Context, s *domain.ApplicationStatus) error {
	return withTx(ctx, r.db, func(tx *sqlx.Tx) error {
		if s.IsDefault {
			if _, err := tx.ExecContext(ctx,
				`UPDATE application_statuses SET is_default = 0 WHERE is_default = 1 AND id <> ?`,
				s.ID); err != nil {
				return err
			}
		} else {
			var n int
			if err := tx.GetContext(ctx, &n,
				`SELECT COUNT(*) FROM application_statuses WHERE is_default = 1 AND id <> ?`,
				s.ID); err != nil {
				return err
			}
			if n == 0 {
				// We're about to demote the only default-flagged row
				// and nothing else is default.
				return ErrStatusNoDefault
			}
		}
		res, err := tx.ExecContext(ctx, `
            UPDATE application_statuses
               SET slug          = ?,
                   name          = ?,
                   kind          = ?,
                   color         = ?,
                   display_order = ?,
                   is_default    = ?
             WHERE id = ?`,
			s.Slug, s.Name, s.Kind, s.Color, s.DisplayOrder, s.IsDefault, s.ID)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// Delete removes a status. Returns ErrStatusInUse if any applications
// or history events still reference it; HR must reassign first.
func (r *StatusRepo) Delete(ctx context.Context, id int64) error {
	var inUse int
	if err := r.db.GetContext(ctx, &inUse, `
        SELECT
            (SELECT COUNT(*) FROM applications WHERE status_id = ?)
          + (SELECT COUNT(*) FROM application_status_events WHERE to_status_id = ? OR from_status_id = ?)
    `, id, id, id); err != nil {
		return err
	}
	if inUse > 0 {
		return fmt.Errorf("%w: %d row(s) reference status %d", ErrStatusInUse, inUse, id)
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM application_statuses WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
