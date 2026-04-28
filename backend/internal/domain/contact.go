package domain

import "time"

// ContactKind identifies the intent of a contact submission.
type ContactKind string

const (
	ContactGeneral     ContactKind = "general"
	ContactPartnership ContactKind = "partnership"
	ContactResearch    ContactKind = "research"
	ContactHiring      ContactKind = "hiring"
)

// ContactSubmission is a message left through the contact form.
type ContactSubmission struct {
	ID        int64       `db:"id"         json:"id"`
	Name      string      `db:"name"       json:"name"`
	Email     string      `db:"email"      json:"email"`
	Company   string      `db:"company"    json:"company,omitempty"`
	Message   string      `db:"message"    json:"message"`
	Kind      ContactKind `db:"kind"       json:"kind"`
	IPAddress string      `db:"ip_address" json:"-"`
	UserAgent string      `db:"user_agent" json:"-"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
}
