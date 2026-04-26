package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// ContactRepo writes contact submissions.
type ContactRepo struct{ db *sqlx.DB }

func NewContactRepo(db *sqlx.DB) *ContactRepo { return &ContactRepo{db: db} }

// Insert persists a new submission and returns the inserted ID.
func (r *ContactRepo) Insert(ctx context.Context, s *domain.ContactSubmission) (int64, error) {
	const q = `
        INSERT INTO contact_submissions
            (name, email, company, message, kind, ip_address, user_agent)
        VALUES
            (:name, :email, :company, :message, :kind, :ip_address, :user_agent)
        RETURNING id
    `
	rows, err := r.db.NamedQueryContext(ctx, q, s)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var id int64
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return id, nil
}
