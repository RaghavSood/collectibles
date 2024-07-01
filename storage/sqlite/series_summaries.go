package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) SeriesSummaries() ([]types.SeriesSummary, error) {
	query := `
WITH item_counts AS (
  SELECT 
    s.slug, 
    COUNT(i.sku) AS item_count
  FROM series s
  JOIN items i ON i.series_slug = s.slug
  GROUP BY s.slug
),
series_tvl AS (
  SELECT 
    s.slug, 
    SUM(op.value) AS tvl
  FROM series s
  JOIN items i ON i.series_slug = s.slug
  JOIN addresses a ON a.sku = i.sku
  JOIN outpoints op ON op.script = a.script
  WHERE op.spending_txid IS NULL
  GROUP BY s.slug
),
item_status AS (
  SELECT 
    s.slug,
    i.sku,
    MAX(CASE WHEN op.txid IS NULL THEN 1 ELSE 0 END) AS unfunded,
    MAX(CASE WHEN op.spending_txid IS NOT NULL THEN 1 ELSE 0 END) AS redeemed,
    MIN(CASE WHEN op.spending_txid IS NULL AND op.txid IS NOT NULL THEN 1 ELSE 0 END) AS unredeemed
  FROM series s
  JOIN items i ON i.series_slug = s.slug
  LEFT JOIN addresses a ON a.sku = i.sku
  LEFT JOIN outpoints op ON op.script = a.script
  GROUP BY s.slug, i.sku
),
status_counts AS (
  SELECT
    slug,
    SUM(unfunded) AS unfunded_count,
    SUM(redeemed) AS redeemed_count,
    SUM(unredeemed) AS unredeemed_count
  FROM item_status
  GROUP BY slug
)
SELECT 
  s.slug, 
  s.name, 
  COALESCE(ic.item_count, 0) AS item_count, 
  COALESCE(st.tvl, 0) AS tvl,
  COALESCE(sc.unfunded_count, 0) AS unfunded,
  COALESCE(sc.redeemed_count, 0) AS redeemed,
  COALESCE(sc.unredeemed_count, 0) AS unredeemed
FROM series s
LEFT JOIN item_counts ic ON s.slug = ic.slug
LEFT JOIN series_tvl st ON s.slug = st.slug
LEFT JOIN status_counts sc ON s.slug = sc.slug;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanSeriesSummaries(rows)
}

func scanSeriesSummaries(rows *sql.Rows) ([]types.SeriesSummary, error) {
	var summaries []types.SeriesSummary
	for rows.Next() {
		var summary types.SeriesSummary
		err := rows.Scan(&summary.Slug, &summary.Name, &summary.ItemCount, &summary.TotalValue, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
