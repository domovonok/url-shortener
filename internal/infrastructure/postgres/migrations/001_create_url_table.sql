-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url (
    id          BIGSERIAL PRIMARY KEY,
    url         TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd
