package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) ItemSummaries() ([]types.ItemSummary, error) {
	query := `WITH item_tvl AS (
  SELECT 
    i.sku, 
    SUM(op.value) AS tvl
  FROM items i
  JOIN addresses a ON a.sku = i.sku
  JOIN outpoints op ON op.script = a.script
  WHERE op.spending_txid IS NULL
  GROUP BY i.sku
),
unspent_outpoints AS (
  SELECT 
    i.sku, 
    COUNT(op.txid) AS unspent
  FROM items i
  JOIN addresses a ON a.sku = i.sku
  JOIN outpoints op ON op.script = a.script
  WHERE op.spending_txid IS NULL
  GROUP BY i.sku
),
spent_outpoints AS (
  SELECT 
    i.sku, 
    COUNT(op.txid) AS spent
  FROM items i
  JOIN addresses a ON a.sku = i.sku
  JOIN outpoints op ON op.script = a.script
  WHERE op.spending_txid IS NOT NULL
  GROUP BY i.sku
),
total_received AS (
  SELECT 
    i.sku, 
    SUM(op.value) AS total_received
  FROM items i
  JOIN addresses a ON a.sku = i.sku
  JOIN outpoints op ON op.script = a.script
  GROUP BY i.sku
),
total_spent AS (
  SELECT 
    i.sku, 
    SUM(op.value) AS total_spent
  FROM items i
  JOIN addresses a ON a.sku = i.sku
  JOIN outpoints op ON op.script = a.script
  WHERE op.spending_txid IS NOT NULL
  GROUP BY i.sku
)
SELECT 
  i.sku, 
  i.serial, 
  s.name AS series_name, 
  s.slug AS series_slug, 
  COALESCE(it.tvl, 0) AS tvl, 
  COALESCE(uo.unspent, 0) AS unspent, 
  COALESCE(so.spent, 0) AS spent, 
  COALESCE(tr.total_received, 0) AS total_received, 
  COALESCE(ts.total_spent, 0) AS total_spent
FROM items i
JOIN series s ON i.series_slug = s.slug
LEFT JOIN item_tvl it ON i.sku = it.sku
LEFT JOIN unspent_outpoints uo ON i.sku = uo.sku
LEFT JOIN spent_outpoints so ON i.sku = so.sku
LEFT JOIN total_received tr ON i.sku = tr.sku
LEFT JOIN total_spent ts ON i.sku = ts.sku;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanItemSummaries(rows)
}

func scanItemSummaries(rows *sql.Rows) ([]types.ItemSummary, error) {
	var summaries []types.ItemSummary
	for rows.Next() {
		var summary types.ItemSummary
		err := rows.Scan(&summary.SKU, &summary.Serial, &summary.SeriesName, &summary.SeriesSlug, &summary.TotalValue, &summary.Unspent, &summary.Spent, &summary.TotalReceived, &summary.TotalSpent)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
