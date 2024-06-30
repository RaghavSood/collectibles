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
)
SELECT 
  c.name, 
  c.slug, 
  COALESCE(sc.series_count, 0) AS series_count, 
  COALESCE(ic.item_count, 0) AS item_count, 
  COALESCE(ct.tvl, 0) AS tvl
FROM creators c
LEFT JOIN series_counts sc ON c.slug = sc.creator_slug
LEFT JOIN item_counts ic ON c.slug = ic.creator_slug
LEFT JOIN creator_tvl ct ON c.slug = ct.creator_slug;`

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
		err := rows.Scan(&summary.Name, &summary.Slug, &summary.SeriesCount, &summary.ItemCount, &summary.TotalValue)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
