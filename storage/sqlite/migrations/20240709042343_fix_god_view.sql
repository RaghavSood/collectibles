-- +goose Up
-- +goose StatementBegin
DROP VIEW god_view;
CREATE VIEW god_view AS
SELECT
  s.name AS series_name,
  s.slug AS series_id,
  json_group_array(DISTINCT(c.name)) AS creators,
  i.sku AS item_id,
  i.serial,
  a.address,
  SUM(op.value) AS total_value,
  MIN(op.block_time) AS first_active,
  MIN(op.spending_block_time) AS redeemed_on,
  SUM(CASE WHEN op.spending_txid IS NULL THEN op.value ELSE 0 END) AS balance
FROM
  series s
JOIN
  series_creators sc ON s.slug = sc.series_slug
JOIN
  creators c ON sc.creator_slug = c.slug
JOIN
  items i ON s.slug = i.series_slug
LEFT JOIN
  addresses a ON i.sku = a.sku
LEFT JOIN
  outpoints op ON a.script = op.script
GROUP BY
  s.name, s.slug, i.sku, i.serial, a.address;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW god_view;
CREATE VIEW god_view AS
SELECT
  s.name AS series_name,
  s.slug AS series_id,
  json_group_array(DISTINCT(c.name)) AS creators,
  i.sku AS item_id,
  i.serial,
  a.address,
  SUM(op.value) AS total_value,
  MIN(op.block_time) AS first_active,
  MAX(op.spending_block_time) AS redeemed_on
FROM
  series s
JOIN
  series_creators sc ON s.slug = sc.series_slug
JOIN
  creators c ON sc.creator_slug = c.slug
JOIN
  items i ON s.slug = i.series_slug
LEFT JOIN
  addresses a ON i.sku = a.sku
LEFT JOIN
  outpoints op ON a.script = op.script
GROUP BY
  s.name, s.slug, i.sku, i.serial, a.address;
-- +goose StatementEnd
