-- +goose Up
-- +goose StatementBegin
ALTER TABLE sessions ADD COLUMN cached_tokens INTEGER NOT NULL DEFAULT 0 CHECK (cached_tokens >= 0);
ALTER TABLE sessions ADD COLUMN total_requests INTEGER NOT NULL DEFAULT 0 CHECK (total_requests >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions DROP COLUMN cached_tokens;
ALTER TABLE sessions DROP COLUMN total_requests;
-- +goose StatementEnd