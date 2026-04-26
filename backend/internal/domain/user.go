package domain

import "time"

// Role enumerates the authorisation roles a user can hold.
type Role string

const (
	// RoleHR can manage job postings via the admin API.
	RoleHR Role = "hr"
)

// User is an authenticated principal allowed to access privileged endpoints.
type User struct {
	ID           int64      `db:"id"            json:"id"`
	Email        string     `db:"email"         json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Role         Role       `db:"role"          json:"role"`
	CreatedAt    time.Time  `db:"created_at"    json:"createdAt"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"lastLoginAt,omitempty"`
}
