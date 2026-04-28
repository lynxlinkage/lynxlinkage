-- +goose Up
-- +goose StatementBegin
CREATE TABLE research_cards (
    id              BIGSERIAL   PRIMARY KEY,
    title           TEXT        NOT NULL,
    summary         TEXT        NOT NULL,
    tags            TEXT        NOT NULL DEFAULT '[]',
    cover_image_url TEXT        NOT NULL DEFAULT '',
    external_url    TEXT        NOT NULL,
    source          TEXT        NOT NULL DEFAULT 'internal',
    published_at    TIMESTAMPTZ NOT NULL,
    display_order   INTEGER     NOT NULL DEFAULT 0
);

CREATE INDEX idx_research_cards_published_at ON research_cards(published_at DESC);
CREATE INDEX idx_research_cards_display_order ON research_cards(display_order);

CREATE TABLE job_postings (
    id                   BIGSERIAL   PRIMARY KEY,
    title                TEXT        NOT NULL,
    team                 TEXT        NOT NULL DEFAULT '',
    location             TEXT        NOT NULL DEFAULT '',
    employment_type      TEXT        NOT NULL DEFAULT 'full_time',
    description_md       TEXT        NOT NULL DEFAULT '',
    apply_url_or_email   TEXT        NOT NULL,
    posted_at            TIMESTAMPTZ NOT NULL,
    is_active            BOOLEAN     NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_job_postings_active_posted ON job_postings(is_active, posted_at DESC);

CREATE TABLE partners (
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT      NOT NULL UNIQUE,
    logo_url      TEXT      NOT NULL,
    website_url   TEXT      NOT NULL DEFAULT '',
    tier          TEXT      NOT NULL DEFAULT 'strategic',
    description   TEXT      NOT NULL DEFAULT '',
    display_order INTEGER   NOT NULL DEFAULT 0
);

CREATE INDEX idx_partners_tier_order ON partners(tier, display_order);

CREATE TABLE contact_submissions (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    company    TEXT        NOT NULL DEFAULT '',
    message    TEXT        NOT NULL,
    kind       TEXT        NOT NULL DEFAULT 'general',
    ip_address TEXT        NOT NULL DEFAULT '',
    user_agent TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_contact_submissions_created_at ON contact_submissions(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contact_submissions;
DROP TABLE IF EXISTS partners;
DROP TABLE IF EXISTS job_postings;
DROP TABLE IF EXISTS research_cards;
-- +goose StatementEnd
