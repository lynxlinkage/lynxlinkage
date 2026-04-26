package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// JobRepo reads job postings.
type JobRepo struct{ db *sqlx.DB }

func NewJobRepo(db *sqlx.DB) *JobRepo { return &JobRepo{db: db} }

// ListActive returns currently active postings, newest first.
func (r *JobRepo) ListActive(ctx context.Context) ([]domain.JobPosting, error) {
	const q = `
        SELECT id, title, team, location, employment_type, description_md,
               apply_url_or_email, posted_at, is_active
        FROM job_postings
        WHERE is_active = 1
        ORDER BY posted_at DESC, id DESC
    `
	out := []domain.JobPosting{}
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single posting by ID, regardless of active status.
func (r *JobRepo) Get(ctx context.Context, id int64) (*domain.JobPosting, error) {
	const q = `
        SELECT id, title, team, location, employment_type, description_md,
               apply_url_or_email, posted_at, is_active
        FROM job_postings WHERE id = ?
    `
	var j domain.JobPosting
	if err := r.db.GetContext(ctx, &j, q, id); err != nil {
		return nil, err
	}
	return &j, nil
}

// Upsert inserts or replaces a job posting. Used by the seed loader.
func (r *JobRepo) Upsert(ctx context.Context, j *domain.JobPosting) error {
	const q = `
        INSERT INTO job_postings
            (id, title, team, location, employment_type, description_md,
             apply_url_or_email, posted_at, is_active)
        VALUES
            (:id, :title, :team, :location, :employment_type, :description_md,
             :apply_url_or_email, :posted_at, :is_active)
        ON CONFLICT(id) DO UPDATE SET
            title=excluded.title,
            team=excluded.team,
            location=excluded.location,
            employment_type=excluded.employment_type,
            description_md=excluded.description_md,
            apply_url_or_email=excluded.apply_url_or_email,
            posted_at=excluded.posted_at,
            is_active=excluded.is_active
    `
	_, err := r.db.NamedExecContext(ctx, q, j)
	return err
}
