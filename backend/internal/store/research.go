package store

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// ResearchRepo reads research cards.
type ResearchRepo struct{ db *sqlx.DB }

func NewResearchRepo(db *sqlx.DB) *ResearchRepo { return &ResearchRepo{db: db} }

// ListOpts narrows a research listing.
type ListResearchOpts struct {
	Tag   string
	Limit int
}

// List returns research cards ordered by display_order, then most recent first.
func (r *ResearchRepo) List(ctx context.Context, opts ListResearchOpts) ([]domain.ResearchCard, error) {
	var (
		args   []any
		where  []string
		query  strings.Builder
	)
	query.WriteString(`SELECT id, title, summary, tags, cover_image_url, external_url,
		source, published_at, display_order FROM research_cards`)

	if opts.Tag != "" {
		// SQLite doesn't have first-class JSON_CONTAINS; we use json_each.
		// Card matches if any element in its tags array equals the filter.
		where = append(where, `EXISTS (SELECT 1 FROM json_each(research_cards.tags) WHERE value = ?)`)
		args = append(args, opts.Tag)
	}
	if len(where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(where, " AND "))
	}
	query.WriteString(` ORDER BY display_order ASC, published_at DESC`)
	if opts.Limit > 0 {
		query.WriteString(` LIMIT ?`)
		args = append(args, opts.Limit)
	}

	out := []domain.ResearchCard{}
	if err := r.db.SelectContext(ctx, &out, query.String(), args...); err != nil {
		return nil, err
	}
	return out, nil
}

// Upsert inserts or replaces a research card. Used by the seed loader.
func (r *ResearchRepo) Upsert(ctx context.Context, c *domain.ResearchCard) error {
	const q = `
        INSERT INTO research_cards
            (id, title, summary, tags, cover_image_url, external_url, source, published_at, display_order)
        VALUES
            (:id, :title, :summary, :tags, :cover_image_url, :external_url, :source, :published_at, :display_order)
        ON CONFLICT(id) DO UPDATE SET
            title=excluded.title,
            summary=excluded.summary,
            tags=excluded.tags,
            cover_image_url=excluded.cover_image_url,
            external_url=excluded.external_url,
            source=excluded.source,
            published_at=excluded.published_at,
            display_order=excluded.display_order
    `
	_, err := r.db.NamedExecContext(ctx, q, c)
	return err
}
