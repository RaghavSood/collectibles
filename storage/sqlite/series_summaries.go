package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) SeriesSummaries() ([]types.SeriesSummary, error) {
	query := `WITH item_counts AS (
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
)
SELECT 
  s.slug, 
  s.name, 
  COALESCE(ic.item_count, 0) AS item_count, 
  COALESCE(st.tvl, 0) AS tvl
FROM series s
LEFT JOIN item_counts ic ON s.slug = ic.slug
LEFT JOIN series_tvl st ON s.slug = st.slug;`

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
		err := rows.Scan(&summary.Slug, &summary.Name, &summary.ItemCount, &summary.TotalValue)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
