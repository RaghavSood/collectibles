-- +goose Up
-- +goose StatementBegin
CREATE TABLE script_queue (
  script TEXT NOT NULL,
  chain TEXT NOT NULL,
  try_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE transaction_queue (
  txid TEXT NOT NULL,
  chain TEXT NOT NULL,
  block_height INTEGER NOT NULL,
  try_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE script_queue;
DROP TABLE transaction_queue;
-- +goose StatementEnd
