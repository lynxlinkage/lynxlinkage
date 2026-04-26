package domain

import "time"

// Application is a single candidate submission against a job posting.
type Application struct {
	ID        int64     `db:"id"          json:"id"`
	JobID     int64     `db:"job_id"      json:"jobId"`
	Name      string    `db:"name"        json:"name"`
	Email     string    `db:"email"       json:"email"`
	Message   string    `db:"message"     json:"message"`
	IPAddress string    `db:"ip_address"  json:"-"`
	UserAgent string    `db:"user_agent"  json:"-"`
	CreatedAt time.Time `db:"created_at"  json:"createdAt"`

	// Files is populated by repo helpers when requested; not stored on
	// the row itself.
	Files []ApplicationFile `db:"-" json:"files,omitempty"`

	// JobTitle is denormalised by joins on list endpoints so the admin
	// UI can render the table without an extra round trip per row.
	JobTitle string `db:"job_title" json:"jobTitle,omitempty"`
}

// ApplicationFile is metadata about an uploaded artifact (CV, cover
// letter, samples, etc.). Bytes live on disk; this row points at them.
type ApplicationFile struct {
	ID            int64     `db:"id"             json:"id"`
	ApplicationID int64     `db:"application_id" json:"applicationId"`
	OriginalName  string    `db:"original_name"  json:"originalName"`
	StoredPath    string    `db:"stored_path"    json:"-"`
	ContentType   string    `db:"content_type"   json:"contentType"`
	SizeBytes     int64     `db:"size_bytes"     json:"sizeBytes"`
	CreatedAt     time.Time `db:"created_at"     json:"createdAt"`
}
