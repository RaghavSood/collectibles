-- +goose Up
-- +goose StatementBegin
CREATE VIEW address_summary AS
SELECT
  a.address,
  a.sku,
  i.series_slug,
  i.serial,
  COUNT(CASE WHEN op.spending_txid IS NULL AND op.txid IS NOT NULL THEN 1 ELSE NULL END) AS unspent,
  COUNT(CASE WHEN op.spending_txid IS NOT NULL THEN 1 ELSE NULL END) AS spent,
  MIN(op.block_time) as first_active,
  MIN(op.spending_block_time) as redeemed_on,
  SUM(CASE WHEN op.spending_txid IS NULL THEN op.value ELSE 0 END) AS total_value
FROM
  addresses a
JOIN
  items i ON a.sku = i.sku
LEFT JOIN
  outpoints op ON a.script = op.script
GROUP BY
  a.address, a.sku, i.series_slug, i.serial;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW address_summary;
-- +goose StatementEnd
