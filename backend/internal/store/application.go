package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// ApplicationRepo persists candidate applications, uploaded artifact
// metadata, and pipeline status changes. The on-disk byte storage is
// owned by package uploads; this repo only deals in rows.
type ApplicationRepo struct{ db *sqlx.DB }

func NewApplicationRepo(db *sqlx.DB) *ApplicationRepo { return &ApplicationRepo{db: db} }

// applicationSelect lists every column the repo reads on List/Get,
// including the joined status fields. Keeping it in one constant
// guarantees List and Get stay in sync.
const applicationSelect = `
    SELECT a.id, a.job_id, a.name, a.email, a.message,
           a.ip_address, a.user_agent, a.created_at,
           a.status_id, a.status_updated_at, a.status_updated_by,
           COALESCE(j.title, '')        AS job_title,
           s.slug                        AS status_slug,
           s.name                        AS status_name,
           s.kind                        AS status_kind,
           s.color                       AS status_color,
           s.display_order               AS status_display_order,
           s.is_default                  AS status_is_default,
           s.created_at                  AS status_created_at
    FROM applications a
    LEFT JOIN job_postings j         ON j.id = a.job_id
    LEFT JOIN application_statuses s ON s.id = a.status_id
`

// applicationRow is what we scan into; the joined status columns are
// nullable because not every row has a status (or the status was
// deleted via ON DELETE SET NULL). We hydrate domain.Application.Status
// from these fields after scanning.
type applicationRow struct {
	domain.Application
	StatusSlug         *string `db:"status_slug"`
	StatusName         *string `db:"status_name"`
	StatusKind         *string `db:"status_kind"`
	StatusColor        *string `db:"status_color"`
	StatusDisplayOrder *int    `db:"status_display_order"`
	StatusIsDefault    *bool   `db:"status_is_default"`
	StatusCreatedAt    *string `db:"status_created_at"`
}

func (r applicationRow) toDomain() domain.Application {
	app := r.Application
	if r.StatusName != nil && app.StatusID != nil {
		s := domain.ApplicationStatus{ID: *app.StatusID, Name: *r.StatusName}
		if r.StatusSlug != nil {
			s.Slug = *r.StatusSlug
		}
		if r.StatusKind != nil {
			s.Kind = domain.ApplicationStatusKind(*r.StatusKind)
		}
		if r.StatusColor != nil {
			s.Color = *r.StatusColor
		}
		if r.StatusDisplayOrder != nil {
			s.DisplayOrder = *r.StatusDisplayOrder
		}
		if r.StatusIsDefault != nil {
			s.IsDefault = *r.StatusIsDefault
		}
		app.Status = &s
	}
	return app
}

// Create inserts an application row and returns the new ID. The caller
// passes in defaultStatusID (the id of the status to apply to brand-new
// rows); if zero, the row is created without a status and shows up as
// "needs triage" in the admin UI.
func (r *ApplicationRepo) Create(ctx context.Context, a *domain.Application, defaultStatusID int64) (int64, error) {
	if defaultStatusID > 0 {
		a.StatusID = &defaultStatusID
	}
	const q = `
        INSERT INTO applications (
            job_id, name, email, message, ip_address, user_agent,
            status_id, status_updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
        RETURNING id
    `
	var id int64
	if err := r.db.QueryRowxContext(ctx, q,
		a.JobID, a.Name, a.Email, a.Message, a.IPAddress, a.UserAgent,
		a.StatusID,
	).Scan(&id); err != nil {
		return 0, err
	}
	a.ID = id

	// Seed the audit trail with the initial status assignment so the
	// detail view shows a complete history.
	if defaultStatusID > 0 {
		if _, err := r.db.ExecContext(ctx, `
            INSERT INTO application_status_events (application_id, from_status_id, to_status_id, actor_id, note)
            VALUES ($1, NULL, $2, NULL, 'submitted')
        `, id, defaultStatusID); err != nil {
			return 0, fmt.Errorf("seed status event: %w", err)
		}
	}
	return id, nil
}

// ExistsRecentApplication returns true when the given email already has
// a submitted application for the same job within the window ending now.
// Email comparison is case-insensitive.
func (r *ApplicationRepo) ExistsRecentApplication(ctx context.Context, jobID int64, email string, window time.Duration) (bool, error) {
	since := time.Now().Add(-window)
	const q = `
        SELECT COUNT(*) FROM applications
        WHERE job_id    = $1
          AND lower(email) = lower($2)
          AND created_at  >= $3
    `
	var n int
	if err := r.db.QueryRowxContext(ctx, q, jobID, email, since).Scan(&n); err != nil {
		return false, err
	}
	return n > 0, nil
}

// AddFile inserts a row pointing at a previously-saved file on disk.
func (r *ApplicationRepo) AddFile(ctx context.Context, f *domain.ApplicationFile) (int64, error) {
	const q = `
        INSERT INTO application_files
            (application_id, original_name, stored_path, content_type, size_bytes)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	var id int64
	if err := r.db.QueryRowxContext(ctx, q,
		f.ApplicationID, f.OriginalName, f.StoredPath, f.ContentType, f.SizeBytes,
	).Scan(&id); err != nil {
		return 0, err
	}
	f.ID = id
	return id, nil
}

// Delete removes an application (and via FK cascade its file rows).
// Bytes on disk must be removed by the caller — this repo doesn't know
// about UPLOAD_DIR.
func (r *ApplicationRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM applications WHERE id = $1`, id)
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

// ListFilter is the bag of optional filters and sort options accepted
// by List. A zero value lists all applications, newest first.
type ListFilter struct {
	JobID    int64
	StatusID int64
	Sort     string // "newest" (default) or "oldest"
	Limit    int
}

// List returns applications matching the filter, joined with the
// related job posting (for job_title) and status (for the badge).
func (r *ApplicationRepo) List(ctx context.Context, f ListFilter) ([]domain.Application, error) {
	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 200
	}

	var (
		where []string
		args  []any
	)
	next := func(v any) string {
		args = append(args, v)
		return "$" + strconv.Itoa(len(args))
	}
	if f.JobID > 0 {
		where = append(where, "a.job_id = "+next(f.JobID))
	}
	if f.StatusID > 0 {
		where = append(where, "a.status_id = "+next(f.StatusID))
	}

	q := applicationSelect
	if len(where) > 0 {
		q += ` WHERE ` + strings.Join(where, " AND ")
	}

	if strings.EqualFold(f.Sort, "oldest") {
		q += ` ORDER BY a.created_at ASC, a.id ASC `
	} else {
		q += ` ORDER BY a.created_at DESC, a.id DESC `
	}
	q += ` LIMIT ` + next(limit)

	rows := []applicationRow{}
	if err := r.db.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, err
	}
	out := make([]domain.Application, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.toDomain())
	}
	return out, nil
}

// Get returns a single application with the joined status. Files and
// history are loaded separately by the caller.
func (r *ApplicationRepo) Get(ctx context.Context, id int64) (*domain.Application, error) {
	q := applicationSelect + ` WHERE a.id = $1`
	var row applicationRow
	if err := r.db.GetContext(ctx, &row, q, id); err != nil {
		return nil, err
	}
	app := row.toDomain()
	return &app, nil
}

// ListFiles returns the files attached to an application.
func (r *ApplicationRepo) ListFiles(ctx context.Context, applicationID int64) ([]domain.ApplicationFile, error) {
	const q = `
        SELECT id, application_id, original_name, stored_path, content_type, size_bytes, created_at
        FROM application_files
        WHERE application_id = $1
        ORDER BY id ASC
    `
	out := []domain.ApplicationFile{}
	if err := r.db.SelectContext(ctx, &out, q, applicationID); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFile returns metadata for a single attachment.
func (r *ApplicationRepo) GetFile(ctx context.Context, id int64) (*domain.ApplicationFile, error) {
	const q = `
        SELECT id, application_id, original_name, stored_path, content_type, size_bytes, created_at
        FROM application_files
        WHERE id = $1
    `
	var f domain.ApplicationFile
	if err := r.db.GetContext(ctx, &f, q, id); err != nil {
		return nil, err
	}
	return &f, nil
}

// ErrSameStatus is returned by UpdateStatus when the requested status
// is already the current one. The handler turns this into a 200 with a
// no-op body so HR's "save" button is idempotent.
var ErrSameStatus = errors.New("store: application already has that status")

// UpdateStatus moves an application to newStatusID and writes one row
// to application_status_events. The whole change is one transaction so
// a row's status and its history can never disagree.
func (r *ApplicationRepo) UpdateStatus(
	ctx context.Context,
	applicationID int64,
	newStatusID int64,
	actorID *int64,
	note string,
) (*domain.ApplicationStatusEvent, error) {
	var event *domain.ApplicationStatusEvent
	err := withTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var current struct {
			StatusID *int64 `db:"status_id"`
		}
		if err := tx.GetContext(ctx, &current,
			`SELECT status_id FROM applications WHERE id = $1`, applicationID); err != nil {
			return err
		}
		if current.StatusID != nil && *current.StatusID == newStatusID {
			return ErrSameStatus
		}

		var exists int
		if err := tx.GetContext(ctx, &exists,
			`SELECT COUNT(*) FROM application_statuses WHERE id = $1`, newStatusID); err != nil {
			return err
		}
		if exists == 0 {
			return ErrNotFound
		}

		if _, err := tx.ExecContext(ctx, `
            UPDATE applications
               SET status_id          = $1,
                   status_updated_at  = CURRENT_TIMESTAMP,
                   status_updated_by  = $2
             WHERE id = $3
        `, newStatusID, actorID, applicationID); err != nil {
			return err
		}

		var eventID int64
		if err := tx.QueryRowxContext(ctx, `
            INSERT INTO application_status_events
                (application_id, from_status_id, to_status_id, actor_id, note)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `, applicationID, current.StatusID, newStatusID, actorID, note).Scan(&eventID); err != nil {
			return err
		}
		event = &domain.ApplicationStatusEvent{
			ID:            eventID,
			ApplicationID: applicationID,
			FromStatusID:  current.StatusID,
			ToStatusID:    newStatusID,
			ActorID:       actorID,
			Note:          note,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return event, nil
}

// ListEvents returns the status-change audit log for one application,
// oldest first, with status names and actor email joined for display.
func (r *ApplicationRepo) ListEvents(ctx context.Context, applicationID int64) ([]domain.ApplicationStatusEvent, error) {
	const q = `
        SELECT e.id, e.application_id, e.from_status_id, e.to_status_id,
               e.actor_id, e.note, e.created_at,
               COALESCE(fs.name, '') AS from_status_name,
               COALESCE(ts.name, '') AS to_status_name,
               COALESCE(u.email, '') AS actor_email
        FROM application_status_events e
        LEFT JOIN application_statuses fs ON fs.id = e.from_status_id
        LEFT JOIN application_statuses ts ON ts.id = e.to_status_id
        LEFT JOIN users u                 ON u.id = e.actor_id
        WHERE e.application_id = $1
        ORDER BY e.created_at ASC, e.id ASC
    `
	out := []domain.ApplicationStatusEvent{}
	if err := r.db.SelectContext(ctx, &out, q, applicationID); err != nil {
		return nil, err
	}
	return out, nil
}
