-- +goose Up
-- +goose StatementBegin
CREATE TABLE kv (
  key TEXT PRIMARY KEY,
  value TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE kv;
-- +goose StatementEnd
