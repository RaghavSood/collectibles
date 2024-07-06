-- +goose Up
-- +goose StatementBegin
DROP VIEW transaction_view;
CREATE VIEW transaction_view AS
SELECT
  op.txid AS txid,
  op.vout AS vout,
  NULL as vin,
  NULL AS original_txid,
  op.value,
  op.block_height,
  op.block_time,
  'incoming' AS transaction_type,
  i.sku,
  i.series_slug,
  i.serial,
  s.name
FROM
  outpoints op
JOIN
  addresses a ON op.script = a.script
JOIN
  items i ON a.sku = i.sku
JOIN
  series s ON i.series_slug = s.slug

UNION ALL

SELECT
  op.spending_txid AS txid,
  NULL AS vout,
  op.spending_vin AS vin,
  op.txid AS original_txid,
  op.value AS value,
  op.spending_block_height AS block_height,
  op.spending_block_time AS block_time,
  'outgoing' AS transaction_type,
  i.sku,
  i.series_slug,
  i.serial,
  s.name
FROM
  outpoints op
JOIN
  addresses a ON op.script = a.script
JOIN
  items i ON a.sku = i.sku
JOIN
  series s ON i.series_slug = s.slug
WHERE
  op.spending_txid IS NOT NULL

ORDER BY
  block_height, block_time;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW transaction_view;
CREATE VIEW transaction_view AS
SELECT
  op.txid AS txid,
  op.vout AS vout,
  NULL as vin,
  NULL AS original_txid,
  op.value,
  op.block_height,
  op.block_time,
  'incoming' AS transaction_type,
  i.sku,
  i.series_slug
FROM
  outpoints op
JOIN
  addresses a ON op.script = a.script
JOIN
  items i ON a.sku = i.sku

UNION ALL

SELECT
  op.spending_txid AS txid,
  NULL AS vout,
  op.spending_vin AS vin,
  op.txid AS original_txid,
  -op.value AS value,
  op.spending_block_height AS block_height,
  op.spending_block_time AS block_time,
  'outgoing' AS transaction_type,
  i.sku,
  i.series_slug
FROM
  outpoints op
JOIN
  addresses a ON op.script = a.script
JOIN
  items i ON a.sku = i.sku
WHERE
  op.spending_txid IS NOT NULL

ORDER BY
  block_height, block_time;
-- +goose StatementEnd
