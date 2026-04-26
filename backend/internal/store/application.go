package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// ApplicationRepo persists candidate applications and uploaded artifact
// metadata. The on-disk byte storage is owned by package uploads; this
// repo only deals in rows.
type ApplicationRepo struct{ db *sqlx.DB }

func NewApplicationRepo(db *sqlx.DB) *ApplicationRepo { return &ApplicationRepo{db: db} }

// Create inserts an application row and returns the new ID. Files are
// inserted separately via AddFile after the bytes are persisted on disk.
func (r *ApplicationRepo) Create(ctx context.Context, a *domain.Application) (int64, error) {
	const q = `
        INSERT INTO applications (job_id, name, email, message, ip_address, user_agent)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	res, err := r.db.ExecContext(ctx, q,
		a.JobID, a.Name, a.Email, a.Message, a.IPAddress, a.UserAgent,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	a.ID = id
	return id, nil
}

// AddFile inserts a row pointing at a previously-saved file on disk.
func (r *ApplicationRepo) AddFile(ctx context.Context, f *domain.ApplicationFile) (int64, error) {
	const q = `
        INSERT INTO application_files
            (application_id, original_name, stored_path, content_type, size_bytes)
        VALUES (?, ?, ?, ?, ?)
    `
	res, err := r.db.ExecContext(ctx, q,
		f.ApplicationID, f.OriginalName, f.StoredPath, f.ContentType, f.SizeBytes,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	f.ID = id
	return id, nil
}

// Delete removes an application (and via FK cascade its file rows).
// Bytes on disk must be removed by the caller — this repo doesn't know
// about UPLOAD_DIR.
func (r *ApplicationRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM applications WHERE id = ?`, id)
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

// List returns applications, newest first, optionally filtered by
// jobID (pass 0 for all). The returned rows are joined to job_postings
// to populate JobTitle so the admin list doesn't need to round-trip.
func (r *ApplicationRepo) List(ctx context.Context, jobID int64, limit int) ([]domain.Application, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	args := []any{}
	q := `
        SELECT a.id, a.job_id, a.name, a.email, a.message,
               a.ip_address, a.user_agent, a.created_at,
               COALESCE(j.title, '') AS job_title
        FROM applications a
        LEFT JOIN job_postings j ON j.id = a.job_id
    `
	if jobID > 0 {
		q += ` WHERE a.job_id = ? `
		args = append(args, jobID)
	}
	q += ` ORDER BY a.created_at DESC, a.id DESC LIMIT ? `
	args = append(args, limit)

	out := []domain.Application{}
	if err := r.db.SelectContext(ctx, &out, q, args...); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single application without its files (use ListFiles
// to populate them).
func (r *ApplicationRepo) Get(ctx context.Context, id int64) (*domain.Application, error) {
	const q = `
        SELECT a.id, a.job_id, a.name, a.email, a.message,
               a.ip_address, a.user_agent, a.created_at,
               COALESCE(j.title, '') AS job_title
        FROM applications a
        LEFT JOIN job_postings j ON j.id = a.job_id
        WHERE a.id = ?
    `
	var a domain.Application
	if err := r.db.GetContext(ctx, &a, q, id); err != nil {
		return nil, err
	}
	return &a, nil
}

// ListFiles returns the files attached to an application.
func (r *ApplicationRepo) ListFiles(ctx context.Context, applicationID int64) ([]domain.ApplicationFile, error) {
	const q = `
        SELECT id, application_id, original_name, stored_path, content_type, size_bytes, created_at
        FROM application_files
        WHERE application_id = ?
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
        WHERE id = ?
    `
	var f domain.ApplicationFile
	if err := r.db.GetContext(ctx, &f, q, id); err != nil {
		return nil, err
	}
	return &f, nil
}
