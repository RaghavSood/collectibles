package sqlite

import (
	"database/sql"
	"time"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) AddressSummariesByItem(sku string) ([]types.AddressSummary, error) {
	rows, err := d.db.Query("SELECT address, sku, series_slug, serial, unspent, spent, first_active, redeemed_on, total_value FROM address_summary WHERE sku = ? ORDER BY serial, first_active, address", sku)
	if err != nil {
		return nil, err
	}

	return scanAddressSummaries(rows)
}

func (d *SqliteBackend) AddressSummariesBySeries(seriesSlug string) ([]types.AddressSummary, error) {
	rows, err := d.db.Query("SELECT address, sku, series_slug, serial, unspent, spent, first_active, redeemed_on, total_value FROM address_summary WHERE series_slug = ? ORDER BY serial, first_active, address", seriesSlug)
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
		var firstActive *string
		var redeemedOn *string
		err := rows.Scan(&summary.Address, &summary.SKU, &summary.SeriesSlug, &summary.Serial, &summary.Unspent, &summary.Spent, &firstActive, &redeemedOn, &summary.TotalValue)
		if err != nil {
			return nil, err
		}

		if firstActive != nil {
			summary.FirstActive, err = time.Parse("2006-01-02 15:04:05-07:00", *firstActive)
			if err != nil {
				return nil, err
			}
		}

		if redeemedOn != nil {
			summary.RedeemedOn, err = time.Parse("2006-01-02 15:04:05-07:00", *redeemedOn)
			if err != nil {
				return nil, err
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
