-- +goose Up
-- +goose StatementBegin
CREATE TABLE outpoint (
  txid TEXT NOT NULL,
  vout INTEGER NOT NULL,
  script TEXT NOT NULL,
  value INTEGER NOT NULL,
  block_height INTEGER NOT NULL,
  block_time DATETIME NOT NULL,
  spending_txid TEXT,
  spending_vin INTEGER,
  spending_block_height INTEGER,
  spending_block_time DATETIME,
  PRIMARY KEY (txid, vout)
  UNIQUE (spending_txid, spending_vin)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outpoint;
-- +goose StatementEnd
