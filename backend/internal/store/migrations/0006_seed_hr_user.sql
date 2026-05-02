-- +goose Up
-- +goose StatementBegin
-- Seed the hr@lynxlinkage.com account. Authentication is delegated to
-- Authelia; the password_hash value is a placeholder that can never match
-- any real password because it is not a valid bcrypt hash.
INSERT INTO users (email, password_hash, role)
VALUES ('hr@lynxlinkage.com', '$AUTHELIA_MANAGED$', 'hr')
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'hr@lynxlinkage.com';
-- +goose StatementEnd
