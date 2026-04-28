package domain

import "time"

// ApplicationStatusKind classifies a status as in-flight (open) or
// terminal (accept / reject). HR is free to add many `open` statuses;
// `accept` and `reject` are the two terminal kinds the system reasons
// about (for badges and pipeline metrics later).
type ApplicationStatusKind string

const (
	StatusKindOpen   ApplicationStatusKind = "open"
	StatusKindAccept ApplicationStatusKind = "accept"
	StatusKindReject ApplicationStatusKind = "reject"
)

// ApplicationStatus is one node in HR's hiring pipeline.
type ApplicationStatus struct {
	ID           int64                 `db:"id"            json:"id"`
	Slug         string                `db:"slug"          json:"slug"`
	Name         string                `db:"name"          json:"name"`
	Kind         ApplicationStatusKind `db:"kind"          json:"kind"`
	Color        string                `db:"color"         json:"color"`
	DisplayOrder int                   `db:"display_order" json:"displayOrder"`
	IsDefault    bool                  `db:"is_default"    json:"isDefault"`
	CreatedAt    time.Time             `db:"created_at"    json:"createdAt"`
}

// ApplicationStatusEvent is a single row in the audit trail of status
// changes for an application.
type ApplicationStatusEvent struct {
	ID            int64     `db:"id"             json:"id"`
	ApplicationID int64     `db:"application_id" json:"applicationId"`
	FromStatusID  *int64    `db:"from_status_id" json:"fromStatusId,omitempty"`
	ToStatusID    int64     `db:"to_status_id"   json:"toStatusId"`
	ActorID       *int64    `db:"actor_id"       json:"actorId,omitempty"`
	Note          string    `db:"note"           json:"note"`
	CreatedAt     time.Time `db:"created_at"     json:"createdAt"`

	// Joined fields for display. None of these are persisted on the
	// event row itself.
	FromStatusName string `db:"from_status_name" json:"fromStatusName,omitempty"`
	ToStatusName   string `db:"to_status_name"   json:"toStatusName,omitempty"`
	ActorEmail     string `db:"actor_email"      json:"actorEmail,omitempty"`
}

// Application is a single candidate submission against a job posting.
type Application struct {
	ID              int64      `db:"id"                 json:"id"`
	JobID           int64      `db:"job_id"             json:"jobId"`
	Name            string     `db:"name"               json:"name"`
	Email           string     `db:"email"              json:"email"`
	Message         string     `db:"message"            json:"message"`
	IPAddress       string     `db:"ip_address"         json:"-"`
	UserAgent       string     `db:"user_agent"         json:"-"`
	CreatedAt       time.Time  `db:"created_at"         json:"createdAt"`
	StatusID        *int64     `db:"status_id"          json:"statusId,omitempty"`
	StatusUpdatedAt *time.Time `db:"status_updated_at"  json:"statusUpdatedAt,omitempty"`
	StatusUpdatedBy *int64     `db:"status_updated_by"  json:"statusUpdatedBy,omitempty"`

	// Status is populated by joins on list/detail endpoints so the
	// admin UI doesn't need an extra round trip to render the badge.
	Status *ApplicationStatus `db:"-" json:"status,omitempty"`

	// Files is populated by repo helpers when requested; not stored on
	// the row itself.
	Files []ApplicationFile `db:"-" json:"files,omitempty"`

	// History is populated by Get on the detail endpoint.
	History []ApplicationStatusEvent `db:"-" json:"history,omitempty"`

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
