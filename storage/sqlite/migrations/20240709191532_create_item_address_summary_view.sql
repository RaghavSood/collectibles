-- +goose Up
-- +goose StatementBegin
CREATE VIEW item_address_summary_view AS
SELECT
  i.sku,
  i.serial,
  CASE
    WHEN COUNT(a.address) = 0 THEN '[]'
    ELSE json_group_array(DISTINCT a.address)
  END AS addresses,
  MIN(op.block_time) AS first_active,
  MIN(op.spending_block_time) AS redeemed_on,
  COALESCE(SUM(CASE WHEN op.spending_txid IS NULL THEN op.value ELSE 0 END), 0) AS total_value,
  s.name AS series_name,
  s.slug AS series_slug
FROM
  items i
JOIN
  series s ON i.series_slug = s.slug
LEFT JOIN
  addresses a ON i.sku = a.sku
LEFT JOIN
  outpoints op ON a.script = op.script
GROUP BY
  i.sku, i.serial, s.name, s.slug;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW item_address_summary_view;
-- +goose StatementEnd
