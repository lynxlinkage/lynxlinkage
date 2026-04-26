-- +goose Up
-- +goose StatementBegin
CREATE TABLE applications (
    id          BIGSERIAL   PRIMARY KEY,
    job_id      BIGINT      NOT NULL REFERENCES job_postings(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    email       TEXT        NOT NULL,
    message     TEXT        NOT NULL DEFAULT '',
    ip_address  TEXT        NOT NULL DEFAULT '',
    user_agent  TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_applications_job_created ON applications(job_id, created_at DESC);
CREATE INDEX idx_applications_created     ON applications(created_at DESC);

-- File metadata. Bytes live on disk under UPLOAD_DIR; the row holds the
-- canonical original filename (for download Content-Disposition) plus
-- the relative storage path the server actually opens.
CREATE TABLE application_files (
    id              BIGSERIAL   PRIMARY KEY,
    application_id  BIGINT      NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    original_name   TEXT        NOT NULL,
    stored_path     TEXT        NOT NULL,
    content_type    TEXT        NOT NULL DEFAULT '',
    size_bytes      BIGINT      NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_application_files_app ON application_files(application_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS application_files;
DROP TABLE IF EXISTS applications;
-- +goose StatementEnd
