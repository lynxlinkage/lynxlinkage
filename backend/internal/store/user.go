package store

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// UserRepo persists authentication users.
type UserRepo struct{ db *sqlx.DB }

func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db: db} }

const userColumns = `id, email, password_hash, role, created_at, last_login_at`

// GetByEmail returns the user with the given (case-insensitive) email or
// sql.ErrNoRows when not found.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT ` + userColumns + ` FROM users WHERE email = $1`
	var u domain.User
	if err := r.db.GetContext(ctx, &u, q, strings.ToLower(strings.TrimSpace(email))); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID returns the user with the given id or sql.ErrNoRows.
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	q := `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	var u domain.User
	if err := r.db.GetContext(ctx, &u, q, id); err != nil {
		return nil, err
	}
	return &u, nil
}

// TouchLogin updates last_login_at to the current UTC time. Errors are
// returned but callers typically log-and-ignore — failing this should
// not block a successful authentication.
func (r *UserRepo) TouchLogin(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET last_login_at = $1 WHERE id = $2`,
		time.Now().UTC(), id)
	return err
}

// Create inserts a new user; the email is lower-cased before storing so
// lookups are case-insensitive.
func (r *UserRepo) Create(ctx context.Context, email, passwordHash string, role domain.Role) (int64, error) {
	const q = `
        INSERT INTO users (email, password_hash, role)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int64
	if err := r.db.QueryRowxContext(ctx, q,
		strings.ToLower(strings.TrimSpace(email)),
		passwordHash,
		string(role),
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// UpdatePassword replaces the bcrypt hash for a user.
func (r *UserRepo) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		passwordHash, id)
	return err
}
