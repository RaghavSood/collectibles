-- +goose Up
-- +goose StatementBegin
CREATE TABLE flags (
  flag_scope TEXT NOT NULL,
  flag_type TEXT NOT NULL,
  flag_key TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_flags_scope_key ON flags (flag_key, flag_scope, flag_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE flags;
-- +goose StatementEnd
