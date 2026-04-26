-- +goose Up
-- +goose StatementBegin
-- HR-customisable hiring workflow. `kind` lets the system distinguish
-- terminal outcomes from in-flight states (open) for badging and stats;
-- HR is free to add as many `open` rows as they want between `unread`
-- and the terminal `accept` / `reject` rows.
CREATE TABLE application_statuses (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    slug          TEXT     NOT NULL UNIQUE,
    name          TEXT     NOT NULL,
    kind          TEXT     NOT NULL DEFAULT 'open',
    color         TEXT     NOT NULL DEFAULT '',
    display_order INTEGER  NOT NULL DEFAULT 0,
    is_default    INTEGER  NOT NULL DEFAULT 0,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_application_statuses_order ON application_statuses(display_order);

-- Sensible starter pipeline. HR can rename, reorder, add or remove via
-- the admin UI. `unread` is the default — exactly one row should carry
-- is_default = 1 at any time; that invariant is enforced in code.
INSERT INTO application_statuses (slug, name, kind, color, display_order, is_default) VALUES
    ('unread',     'Unread',     'open',   '#64748b',  0, 1),
    ('reviewing',  'Reviewing',  'open',   '#0ea5e9', 10, 0),
    ('shortlist',  'Shortlist',  'open',   '#6366f1', 20, 0),
    ('interview',  'Interview',  'open',   '#f59e0b', 30, 0),
    ('offer',      'Offer',      'open',   '#8b5cf6', 40, 0),
    ('accepted',   'Accepted',   'accept', '#16a34a', 50, 0),
    ('rejected',   'Rejected',   'reject', '#dc2626', 60, 0);

-- Application columns. Nullable so ON DELETE SET NULL on the FK works
-- without losing the row; the API treats NULL as "needs triage" and
-- the create handler always sets it to the current default at insert
-- time so live rows always have a status.
ALTER TABLE applications ADD COLUMN status_id INTEGER REFERENCES application_statuses(id) ON DELETE SET NULL;
ALTER TABLE applications ADD COLUMN status_updated_at DATETIME;
ALTER TABLE applications ADD COLUMN status_updated_by INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Backfill: every existing row starts as the default ("unread").
UPDATE applications
SET status_id          = (SELECT id FROM application_statuses WHERE slug = 'unread'),
    status_updated_at  = created_at;

CREATE INDEX idx_applications_status ON applications(status_id);

-- Audit trail. Each status change writes one row so HR can see how a
-- candidate moved through the pipeline (and who moved them).
CREATE TABLE application_status_events (
    id              INTEGER  PRIMARY KEY AUTOINCREMENT,
    application_id  INTEGER  NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    from_status_id  INTEGER           REFERENCES application_statuses(id) ON DELETE SET NULL,
    to_status_id    INTEGER  NOT NULL REFERENCES application_statuses(id) ON DELETE CASCADE,
    actor_id        INTEGER           REFERENCES users(id) ON DELETE SET NULL,
    note            TEXT     NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_app_status_events_app ON application_status_events(application_id, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS application_status_events;
DROP INDEX IF EXISTS idx_applications_status;
-- SQLite cannot drop columns; rolling back this migration requires a full
-- table rewrite. Document and leave columns in place if downgrading.
DROP TABLE IF EXISTS application_statuses;
-- +goose StatementEnd
