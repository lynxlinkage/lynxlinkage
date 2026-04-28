package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
)

// PartnerRepo reads partners.
type PartnerRepo struct{ db *sqlx.DB }

func NewPartnerRepo(db *sqlx.DB) *PartnerRepo { return &PartnerRepo{db: db} }

// List returns all partners, ordered by tier then display_order.
func (r *PartnerRepo) List(ctx context.Context) ([]domain.Partner, error) {
	const q = `
        SELECT id, name, logo_url, website_url, tier, description, display_order
        FROM partners
        ORDER BY
            CASE tier
                WHEN 'strategic' THEN 1
                WHEN 'exchange'  THEN 2
                WHEN 'broker'    THEN 3
                WHEN 'tech'      THEN 4
                ELSE 5
            END,
            display_order ASC,
            name ASC
    `
	out := []domain.Partner{}
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

// Upsert inserts or replaces a partner by unique name. Used by the seed loader.
func (r *PartnerRepo) Upsert(ctx context.Context, p *domain.Partner) error {
	const q = `
        INSERT INTO partners
            (name, logo_url, website_url, tier, description, display_order)
        VALUES
            (:name, :logo_url, :website_url, :tier, :description, :display_order)
        ON CONFLICT(name) DO UPDATE SET
            logo_url=excluded.logo_url,
            website_url=excluded.website_url,
            tier=excluded.tier,
            description=excluded.description,
            display_order=excluded.display_order
    `
	_, err := r.db.NamedExecContext(ctx, q, p)
	return err
}

// DeleteNotInNames removes rows whose name is not in keep. If keep is empty, all
// partners are removed. Used by seed so the DB matches the YAML set.
func (r *PartnerRepo) DeleteNotInNames(ctx context.Context, keep []string) (int64, error) {
	if len(keep) == 0 {
		res, err := r.db.ExecContext(ctx, `DELETE FROM partners`)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}
	ph := make([]string, len(keep))
	args := make([]any, len(keep))
	for i, name := range keep {
		ph[i] = fmt.Sprintf("$%d", i+1)
		args[i] = name
	}
	q := `DELETE FROM partners WHERE name NOT IN (` + strings.Join(ph, ", ") + `)`
	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
