-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
    id          BIGSERIAL PRIMARY KEY,
    url         TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd