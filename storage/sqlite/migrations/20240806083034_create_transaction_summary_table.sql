-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_summary_c (
  txid TEXT NOT NULL,
  vout INTEGER,
  vin INTEGER,
  original_txid TEXT,
  value INTEGER NOT NULL,
  block_height INTEGER NOT NULL,
  block_time DATETIME NOT NULL,
  transaction_type TEXT NOT NULL,
  sku TEXT NOT NULL,
  series_slug TEXT NOT NULL,
  serial TEXT,
  name TEXT NOT NULL
);

CREATE INDEX transaction_summary_c_sku ON transaction_summary_c(sku);
CREATE INDEX transaction_summary_c_series_slug ON transaction_summary_c(series_slug);
CREATE INDEX transaction_summary_c_block_height ON transaction_summary_c(block_height);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction_summary_c;
-- +goose StatementEnd
