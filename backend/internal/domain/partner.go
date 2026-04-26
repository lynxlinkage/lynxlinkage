package domain

// PartnerTier groups partners on the partners page.
type PartnerTier string

const (
	TierStrategic PartnerTier = "strategic"
	TierExchange  PartnerTier = "exchange"
	TierBroker    PartnerTier = "broker"
	TierTech      PartnerTier = "tech"
)

// Partner represents a logo + link entry on the partners page.
type Partner struct {
	ID           int64       `db:"id"            json:"id"`
	Name         string      `db:"name"          json:"name"`
	LogoURL      string      `db:"logo_url"      json:"logoUrl"`
	WebsiteURL   string      `db:"website_url"   json:"websiteUrl,omitempty"`
	Tier         PartnerTier `db:"tier"          json:"tier"`
	Description  string      `db:"description"   json:"description,omitempty"`
	DisplayOrder int         `db:"display_order" json:"displayOrder"`
}
