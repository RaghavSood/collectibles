package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) CreatorSummaries() ([]types.CreatorSummary, error) {
	query := `WITH series_counts AS (
  SELECT 
    sc.creator_slug, 
    COUNT(DISTINCT sc.series_slug) AS series_count
  FROM series_creators sc
  GROUP BY sc.creator_slug
),
item_counts AS (
  SELECT 
    sc.creator_slug, 
    COUNT(i.sku) AS item_count
  FROM series_creators sc
  JOIN series s ON s.slug = sc.series_slug
  JOIN items i ON i.series_slug = s.slug
  GROUP BY sc.creator_slug
),
creator_tvl AS (
  SELECT 
    sc.creator_slug, 
    SUM(op.value) AS tvl
  FROM series_creators sc
  JOIN series s ON s.slug = sc.series_slug
  JOIN items i ON i.series_slug = s.slug
  JOIN addresses a ON a.sku = i.sku
  JOIN outpoints op ON op.script = a.script
  WHERE op.spending_txid IS NULL
  GROUP BY sc.creator_slug
),
item_status AS (
  SELECT 
    sc.creator_slug,
    i.sku,
    MAX(CASE WHEN op.txid IS NULL THEN 1 ELSE 0 END) AS unfunded,
    MAX(CASE WHEN op.spending_txid IS NOT NULL THEN 1 ELSE 0 END) AS redeemed,
    MIN(CASE WHEN op.spending_txid IS NULL AND op.txid IS NOT NULL THEN 1 ELSE 0 END) AS unredeemed
  FROM series_creators sc
  JOIN series s ON s.slug = sc.series_slug
  JOIN items i ON i.series_slug = s.slug
  LEFT JOIN addresses a ON a.sku = i.sku
  LEFT JOIN outpoints op ON op.script = a.script
  GROUP BY sc.creator_slug, i.sku
),
status_counts AS (
  SELECT
    creator_slug,
    SUM(unfunded) AS unfunded_count,
    SUM(redeemed) AS redeemed_count,
    SUM(unredeemed) AS unredeemed_count
  FROM item_status
  GROUP BY creator_slug
)
SELECT 
  c.name, 
  c.slug, 
  COALESCE(sc.series_count, 0) AS series_count, 
  COALESCE(ic.item_count, 0) AS item_count, 
  COALESCE(ct.tvl, 0) AS tvl,
  COALESCE(stc.unfunded_count, 0) AS unfunded,
  COALESCE(stc.redeemed_count, 0) AS redeemed,
  COALESCE(stc.unredeemed_count, 0) AS unredeemed
FROM creators c
LEFT JOIN series_counts sc ON c.slug = sc.creator_slug
LEFT JOIN item_counts ic ON c.slug = ic.creator_slug
LEFT JOIN creator_tvl ct ON c.slug = ct.creator_slug
LEFT JOIN status_counts stc ON c.slug = stc.creator_slug;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanCreatorSummaries(rows)
}

func scanCreatorSummaries(rows *sql.Rows) ([]types.CreatorSummary, error) {
	var summaries []types.CreatorSummary
	for rows.Next() {
		var summary types.CreatorSummary
		err := rows.Scan(&summary.Name, &summary.Slug, &summary.SeriesCount, &summary.ItemCount, &summary.TotalValue, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
