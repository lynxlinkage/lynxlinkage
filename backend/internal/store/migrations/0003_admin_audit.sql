-- +goose Up
-- +goose StatementBegin

-- Track creation, last edit, and the HR user responsible for each.
-- SQLite's ALTER TABLE ADD COLUMN can't use CURRENT_TIMESTAMP as a
-- default (non-constant defaults are disallowed), so columns are
-- nullable; the Go layer always populates them on writes and we backfill
-- existing rows below.
ALTER TABLE job_postings ADD COLUMN created_at DATETIME;
ALTER TABLE job_postings ADD COLUMN updated_at DATETIME;
ALTER TABLE job_postings ADD COLUMN created_by INTEGER REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE job_postings ADD COLUMN updated_by INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Backfill timestamps for rows seeded before this migration.
UPDATE job_postings
SET created_at = COALESCE(created_at, posted_at),
    updated_at = COALESCE(updated_at, posted_at);

CREATE INDEX idx_job_postings_updated_at ON job_postings(updated_at DESC);

-- Login activity: useful for spotting unused accounts and for audit.
ALTER TABLE users ADD COLUMN last_login_at DATETIME;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SQLite cannot drop columns without a full table rebuild; the down-path
-- is a no-op. Roll forward, not back.
-- +goose StatementEnd
