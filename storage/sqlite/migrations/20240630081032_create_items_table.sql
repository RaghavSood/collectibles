-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
  sku TEXT NOT NULL UNIQUE,
  series_slug TEXT NOT NULL,
  serial TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (series_slug) REFERENCES series (slug) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd
