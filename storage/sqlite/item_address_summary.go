package sqlite

import (
	"database/sql"
	"time"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) ItemAddressSummariesBySeries(seriesSlug string) ([]types.ItemAddressSummary, error) {
	query := `SELECT sku, serial, addresses, first_active, redeemed_on, total_value, series_name, series_slug FROM item_address_summary_c WHERE series_slug = ?;`

	rows, err := d.db.Query(query, seriesSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItemAddressSummary(rows)
}

func (d *SqliteBackend) ItemAddressSummariesByRedeemedOn(redeemedOn time.Time) ([]types.ItemAddressSummary, error) {
	query := `SELECT sku, serial, addresses, first_active, redeemed_on, total_value, series_name, series_slug FROM item_address_summary_c WHERE redeemed_on = ?;`

	rows, err := d.db.Query(query, redeemedOn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItemAddressSummary(rows)
}

func scanItemAddressSummary(rows *sql.Rows) ([]types.ItemAddressSummary, error) {
	var ias []types.ItemAddressSummary
	for rows.Next() {
		var ia types.ItemAddressSummary
		var firstActive *string
		var redeemedOn *string
		err := rows.Scan(&ia.SKU, &ia.Serial, &ia.Addresses, &firstActive, &redeemedOn, &ia.TotalValue, &ia.SeriesName, &ia.SeriesSlug)
		if err != nil {
			return nil, err
		}

		if firstActive != nil {
			ia.FirstActive, err = time.Parse("2006-01-02T15:04:05-07:00", *firstActive)
			if err != nil {
				return nil, err
			}
		}

		if redeemedOn != nil {
			ia.RedeemedOn, err = time.Parse("2006-01-02T15:04:05-07:00", *redeemedOn)
			if err != nil {
				return nil, err
			}
		}
		ias = append(ias, ia)
	}

	return ias, nil
}
