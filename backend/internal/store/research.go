package store

import (
	"context"
	"strconv"
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
		args  []any
		where []string
		query strings.Builder
	)
	query.WriteString(`SELECT id, title, summary, tags, cover_image_url, external_url,
		source, published_at, display_order FROM research_cards`)

	next := func(v any) string {
		args = append(args, v)
		return "$" + strconv.Itoa(len(args))
	}

	if opts.Tag != "" {
		// tags is stored as a JSON-encoded TEXT array; cast at query
		// time to test membership without changing the column type.
		where = append(where,
			`EXISTS (SELECT 1 FROM jsonb_array_elements_text(research_cards.tags::jsonb) AS v WHERE v = `+next(opts.Tag)+`)`)
	}
	if len(where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(where, " AND "))
	}
	query.WriteString(` ORDER BY display_order ASC, published_at DESC`)
	if opts.Limit > 0 {
		query.WriteString(` LIMIT ` + next(opts.Limit))
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
	if _, err := r.db.NamedExecContext(ctx, q, c); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `
        SELECT setval(
            pg_get_serial_sequence('research_cards', 'id'),
            COALESCE((SELECT MAX(id) FROM research_cards), 1),
            true)
    `)
	return err
}
