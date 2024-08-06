package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) AddressSummariesByItem(sku string) ([]types.AddressSummary, error) {
	rows, err := d.db.Query("SELECT address, sku, series_slug, serial, unspent, spent, first_active, redeemed_on, total_value FROM address_summary_c WHERE sku = ? ORDER BY serial, first_active, address", sku)
	if err != nil {
		return nil, err
	}

	return scanAddressSummaries(rows)
}

func (d *SqliteBackend) AddressSummariesBySeries(seriesSlug string) ([]types.AddressSummary, error) {
	rows, err := d.db.Query("SELECT address, sku, series_slug, serial, unspent, spent, first_active, redeemed_on, total_value FROM address_summary_c WHERE series_slug = ? ORDER BY serial, first_active, address", seriesSlug)
	if err != nil {
		return nil, err
	}

	return scanAddressSummaries(rows)
}

func scanAddressSummaries(rows *sql.Rows) ([]types.AddressSummary, error) {
	defer rows.Close()

	var summaries []types.AddressSummary

	for rows.Next() {
		var summary types.AddressSummary
		err := rows.Scan(&summary.Address, &summary.SKU, &summary.SeriesSlug, &summary.Serial, &summary.Unspent, &summary.Spent, &summary.FirstActive, &summary.RedeemedOn, &summary.TotalValue)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
