package store

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// JobRepo reads and writes job postings.
type JobRepo struct{ db *sqlx.DB }

func NewJobRepo(db *sqlx.DB) *JobRepo { return &JobRepo{db: db} }

// jobColumns is the ordered SELECT column list shared by all readers.
const jobColumns = `id, title, team, location, employment_type, description_md,
       apply_url_or_email, posted_at, is_active,
       created_at, updated_at, created_by, updated_by`

// ListActive returns currently active postings, newest first.
func (r *JobRepo) ListActive(ctx context.Context) ([]domain.JobPosting, error) {
	q := `SELECT ` + jobColumns + `
        FROM job_postings
        WHERE is_active = TRUE
        ORDER BY posted_at DESC, id DESC`
	out := []domain.JobPosting{}
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

// ListAll returns every posting (active and inactive), most recently
// edited first. Used by the admin UI.
func (r *JobRepo) ListAll(ctx context.Context) ([]domain.JobPosting, error) {
	q := `SELECT ` + jobColumns + `
        FROM job_postings
        ORDER BY is_active DESC,
                 COALESCE(updated_at, posted_at) DESC,
                 id DESC`
	out := []domain.JobPosting{}
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

// Create inserts a posting and returns the new ID. actorID is the user
// performing the action; pass nil for system writes (e.g. seeding).
func (r *JobRepo) Create(ctx context.Context, j *domain.JobPosting, actorID *int64) (int64, error) {
	now := time.Now().UTC()
	const q = `
        INSERT INTO job_postings
            (title, team, location, employment_type, description_md,
             apply_url_or_email, posted_at, is_active,
             created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `
	var id int64
	if err := r.db.QueryRowxContext(ctx, q,
		j.Title, j.Team, j.Location, string(j.EmploymentType),
		j.DescriptionMD, j.ApplyURLOrEmail, j.PostedAt, j.IsActive,
		now, now, actorID, actorID,
	).Scan(&id); err != nil {
		return 0, err
	}
	j.ID = id
	j.CreatedAt = &now
	j.UpdatedAt = &now
	j.CreatedBy = actorID
	j.UpdatedBy = actorID
	return id, nil
}

// Delete removes a posting permanently. Returns ErrNotFound when the id
// does not exist.
func (r *JobRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM job_postings WHERE id = $1`, id)
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

// Update replaces the contents of an existing posting. Returns
// ErrNotFound if the id does not exist. actorID is the user performing
// the edit; pass nil for system writes.
func (r *JobRepo) Update(ctx context.Context, j *domain.JobPosting, actorID *int64) error {
	now := time.Now().UTC()
	const q = `
        UPDATE job_postings SET
            title = $1, team = $2, location = $3, employment_type = $4,
            description_md = $5, apply_url_or_email = $6,
            is_active = $7,
            updated_at = $8, updated_by = $9
        WHERE id = $10
    `
	res, err := r.db.ExecContext(ctx, q,
		j.Title, j.Team, j.Location, string(j.EmploymentType),
		j.DescriptionMD, j.ApplyURLOrEmail, j.IsActive,
		now, actorID,
		j.ID,
	)
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
	j.UpdatedAt = &now
	j.UpdatedBy = actorID
	return nil
}

// Get returns a single posting by ID, regardless of active status.
func (r *JobRepo) Get(ctx context.Context, id int64) (*domain.JobPosting, error) {
	q := `SELECT ` + jobColumns + ` FROM job_postings WHERE id = $1`
	var j domain.JobPosting
	if err := r.db.GetContext(ctx, &j, q, id); err != nil {
		return nil, err
	}
	return &j, nil
}

// Upsert inserts or replaces a job posting. Used by the seed loader.
// Audit columns are populated from the row's PostedAt as a best-effort
// timestamp; created_by / updated_by remain NULL for seeded rows.
func (r *JobRepo) Upsert(ctx context.Context, j *domain.JobPosting) error {
	const q = `
        INSERT INTO job_postings
            (id, title, team, location, employment_type, description_md,
             apply_url_or_email, posted_at, is_active,
             created_at, updated_at)
        VALUES
            (:id, :title, :team, :location, :employment_type, :description_md,
             :apply_url_or_email, :posted_at, :is_active,
             :posted_at, :posted_at)
        ON CONFLICT(id) DO UPDATE SET
            title=excluded.title,
            team=excluded.team,
            location=excluded.location,
            employment_type=excluded.employment_type,
            description_md=excluded.description_md,
            apply_url_or_email=excluded.apply_url_or_email,
            posted_at=excluded.posted_at,
            is_active=excluded.is_active,
            updated_at=excluded.posted_at
    `
	if _, err := r.db.NamedExecContext(ctx, q, j); err != nil {
		return err
	}
	// BIGSERIAL sequences are not advanced by manual id inserts, so
	// realign the sequence after seeding to avoid duplicate-key errors
	// the next time a posting is created via Create.
	_, err := r.db.ExecContext(ctx, `
        SELECT setval(
            pg_get_serial_sequence('job_postings', 'id'),
            COALESCE((SELECT MAX(id) FROM job_postings), 1),
            true)
    `)
	return err
}
