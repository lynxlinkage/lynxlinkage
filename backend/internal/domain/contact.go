package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ContactKind identifies the intent of a contact submission.
type ContactKind string

const (
	ContactGeneral     ContactKind = "general"
	ContactPartnership ContactKind = "partnership"
	ContactResearch    ContactKind = "research"
	ContactHiring      ContactKind = "hiring"
	ContactMFTIK       ContactKind = "mftik"
)

const (
	mftikSlugMinLen = 3
	mftikSlugMaxLen = 32
)

// mftikSlugPattern matches a lowercase DNS-label-style instance slug:
// starts/ends with [a-z0-9], optional interior hyphens, length checked separately.
var mftikSlugPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// mftikHostInMessage finds `{slug}.mftik.app` mentions in free-text messages.
// Case-insensitive; requires word-ish boundaries so we do not match
// substrings of longer hostnames.
var mftikHostInMessage = regexp.MustCompile(`(?i)(?:^|[^a-z0-9-])([a-z0-9](?:[a-z0-9-]*[a-z0-9])?)\.mftik\.app(?:[^a-z0-9-]|$)`)

// reservedMFTIKSlugs are instance names that must never appear claimable via
// the contact / request-access flow. This is inbox validation only — not
// provisioning or uniqueness enforcement.
var reservedMFTIKSlugs = map[string]struct{}{
	"www": {}, "admin": {}, "api": {}, "root": {}, "mail": {}, "docs": {},
	"mftik": {}, "ftp": {}, "status": {}, "support": {}, "help": {},
	"blog": {}, "cdn": {}, "static": {},
}

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

// ValidateContactMFTIKFields enforces MFTIK instance-slug rules for contact
// submissions. Call when kind is mftik or when an optional slug field is set.
// Empty slug is allowed (omitted); invalid/reserved slugs return a clear error
// suitable for an HTTP 400 body. Does not persist or provision anything.
func ValidateContactMFTIKFields(kind ContactKind, slug, message string) error {
	slug = strings.TrimSpace(slug)
	if slug != "" || kind == ContactMFTIK {
		if err := ValidateMFTIKSlug(slug); err != nil {
			return err
		}
	}
	if kind == ContactMFTIK {
		if reserved := ReservedMFTIKHostInMessage(message); reserved != "" {
			return errors.New("slug is reserved")
		}
	}
	return nil
}

// ValidateMFTIKSlug checks an optional MFTIK instance slug.
// An empty slug is treated as omitted and always valid.
func ValidateMFTIKSlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil
	}
	n := utf8.RuneCountInString(slug)
	if n < mftikSlugMinLen || n > mftikSlugMaxLen {
		return errors.New("slug must be between 3 and 32 characters")
	}
	if !mftikSlugPattern.MatchString(slug) {
		return errors.New("slug must be lowercase letters, digits, and hyphens, and must not start or end with a hyphen")
	}
	if IsReservedMFTIKSlug(slug) {
		return errors.New("slug is reserved")
	}
	return nil
}

// IsReservedMFTIKSlug reports whether slug is in the reserved set.
func IsReservedMFTIKSlug(slug string) bool {
	_, ok := reservedMFTIKSlugs[strings.ToLower(strings.TrimSpace(slug))]
	return ok
}

// ReservedMFTIKHostInMessage returns a reserved slug if the message clearly
// requests a reserved host like `www.mftik.app`. Returns "" when none found.
// Only reserved hosts are flagged so ordinary contact text is not rejected
// for mentioning a hypothetical or valid hostname.
func ReservedMFTIKHostInMessage(message string) string {
	matches := mftikHostInMessage.FindAllStringSubmatch(message, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		slug := strings.ToLower(m[1])
		if IsReservedMFTIKSlug(slug) {
			return slug
		}
	}
	return ""
}
