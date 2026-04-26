-- +goose Up
-- +goose StatementBegin

-- Audit columns capturing creation, last edit, and the HR user
-- responsible for each posting. Columns are nullable because rows
-- inserted before this migration (e.g. via the seed loader) have no
-- known editor; the Go layer always populates them on writes and we
-- backfill timestamps below.
ALTER TABLE job_postings ADD COLUMN created_at TIMESTAMPTZ;
ALTER TABLE job_postings ADD COLUMN updated_at TIMESTAMPTZ;
ALTER TABLE job_postings ADD COLUMN created_by BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE job_postings ADD COLUMN updated_by BIGINT REFERENCES users(id) ON DELETE SET NULL;

UPDATE job_postings
SET created_at = COALESCE(created_at, posted_at),
    updated_at = COALESCE(updated_at, posted_at);

CREATE INDEX idx_job_postings_updated_at ON job_postings(updated_at DESC);

-- Login activity: useful for spotting unused accounts and for audit.
ALTER TABLE users ADD COLUMN last_login_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS last_login_at;
DROP INDEX IF EXISTS idx_job_postings_updated_at;
ALTER TABLE job_postings DROP COLUMN IF EXISTS updated_by;
ALTER TABLE job_postings DROP COLUMN IF EXISTS created_by;
ALTER TABLE job_postings DROP COLUMN IF EXISTS updated_at;
ALTER TABLE job_postings DROP COLUMN IF EXISTS created_at;
-- +goose StatementEnd
