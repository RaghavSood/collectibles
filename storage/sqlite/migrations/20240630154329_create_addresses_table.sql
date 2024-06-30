-- +goose Up
-- +goose StatementBegin
CREATE TABLE addresses (
  sku TEXT NOT NULL,
  script TEXT NOT NULL,
  address TEXT NOT NULL,
  chain TEXT NOT NULL,
  fast_block_height INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY (sku) REFERENCES items(sku) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE addresses;
-- +goose StatementEnd
