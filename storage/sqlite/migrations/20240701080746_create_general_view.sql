-- +goose Up
-- +goose StatementBegin
CREATE VIEW general_statistics AS
SELECT
  (SELECT COUNT(*) FROM creators) AS creators,
  (SELECT COUNT(*) FROM series) AS series,
  (SELECT COUNT(*) FROM items) AS items,
  (SELECT COUNT(*) FROM addresses) AS addresses,
  (SELECT SUM(value) FROM outpoints WHERE spending_txid IS NULL) AS total_value,
  (SELECT SUM(value) FROM outpoints WHERE spending_txid IS NOT NULL) AS total_redeemed;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW general_statistics;
-- +goose StatementEnd
