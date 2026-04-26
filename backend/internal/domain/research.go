package domain

import "time"

// ResearchSource categorises where a research card ultimately links to.
type ResearchSource string

const (
	SourceMedium   ResearchSource = "medium"
	SourceInternal ResearchSource = "internal"
	SourceExternal ResearchSource = "external"
)

// ResearchCard represents a public research summary card surfaced on the
// landing page. Each card links to the full article hosted elsewhere
// (Medium, the firm's research platform, etc.).
type ResearchCard struct {
	ID            int64          `db:"id"             json:"id"`
	Title         string         `db:"title"          json:"title"`
	Summary       string         `db:"summary"        json:"summary"`
	Tags          StringSlice    `db:"tags"           json:"tags"`
	CoverImageURL string         `db:"cover_image_url" json:"coverImageUrl,omitempty"`
	ExternalURL   string         `db:"external_url"   json:"externalUrl"`
	Source        ResearchSource `db:"source"         json:"source"`
	PublishedAt   time.Time      `db:"published_at"   json:"publishedAt"`
	DisplayOrder  int            `db:"display_order"  json:"displayOrder"`
}
