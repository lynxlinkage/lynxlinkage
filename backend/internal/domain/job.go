package domain

import "time"

// EmploymentType is the work arrangement for a job posting.
type EmploymentType string

const (
	EmploymentFullTime   EmploymentType = "full_time"
	EmploymentPartTime   EmploymentType = "part_time"
	EmploymentContract   EmploymentType = "contract"
	EmploymentInternship EmploymentType = "internship"
)

// JobPosting represents an open role on the hiring page.
type JobPosting struct {
	ID               int64          `db:"id"                 json:"id"`
	Title            string         `db:"title"              json:"title"`
	Team             string         `db:"team"               json:"team"`
	Location         string         `db:"location"           json:"location"`
	EmploymentType   EmploymentType `db:"employment_type"    json:"employmentType"`
	DescriptionMD    string         `db:"description_md"     json:"descriptionMd"`
	ApplyURLOrEmail  string         `db:"apply_url_or_email" json:"applyUrlOrEmail"`
	PostedAt         time.Time      `db:"posted_at"          json:"postedAt"`
	IsActive         bool           `db:"is_active"          json:"isActive"`
}
