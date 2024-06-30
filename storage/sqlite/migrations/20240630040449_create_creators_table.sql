-- +goose Up
-- +goose StatementBegin
CREATE TABLE creators (
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE creators;
-- +goose StatementEnd
