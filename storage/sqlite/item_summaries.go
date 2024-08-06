package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) ItemSummaries() ([]types.ItemSummary, error) {
	query := `SELECT sku, serial, series_name, series_slug, tvl, unspent, spent, total_received, total_spent, unfunded, redeemed, unredeemed FROM item_summary_c;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanItemSummaries(rows)
}

func (d *SqliteBackend) ItemSummariesBySeries(slug string) ([]types.ItemSummary, error) {
	query := `SELECT sku, serial, series_name, series_slug, tvl, unspent, spent, total_received, total_spent, unfunded, redeemed, unredeemed FROM item_summary_c WHERE series_slug = ?;`

	rows, err := d.db.Query(query, slug)
	if err != nil {
		return nil, err
	}

	return scanItemSummaries(rows)
}

func (d *SqliteBackend) ItemSummary(sku string) (*types.ItemSummary, error) {
	query := `SELECT sku, serial, series_name, series_slug, tvl, unspent, spent, total_received, total_spent, unfunded, redeemed, unredeemed FROM item_summary_c WHERE sku = ?;`

	row := d.db.QueryRow(query, sku)
	var summary types.ItemSummary
	err := row.Scan(&summary.SKU, &summary.Serial, &summary.SeriesName, &summary.SeriesSlug, &summary.TotalValue, &summary.Unspent, &summary.Spent, &summary.TotalReceived, &summary.TotalSpent, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func scanItemSummaries(rows *sql.Rows) ([]types.ItemSummary, error) {
	var summaries []types.ItemSummary
	for rows.Next() {
		var summary types.ItemSummary
		err := rows.Scan(&summary.SKU, &summary.Serial, &summary.SeriesName, &summary.SeriesSlug, &summary.TotalValue, &summary.Unspent, &summary.Spent, &summary.TotalReceived, &summary.TotalSpent, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
