package sqlite

import (
	"fmt"
	"time"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) GodView() ([]types.GodView, error) {
	rows, err := d.db.Query(`
		SELECT series_name, series_id, creators, item_id, serial, address, total_value, first_active, redeemed_on
		FROM god_view
		ORDER BY series_id, item_id, serial;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query god view: %w", err)
	}
	defer rows.Close()

	var godView []types.GodView
	for rows.Next() {
		var gv types.GodView
		var firstActive *string
		var redeemedOn *string
		err := rows.Scan(&gv.SeriesName, &gv.SeriesID, &gv.Creators, &gv.ItemID, &gv.Serial, &gv.Address, &gv.TotalValue, &firstActive, &redeemedOn)
		if err != nil {
			return nil, fmt.Errorf("failed to scan god view: %w", err)
		}

		if firstActive != nil {
			gv.FirstActive, err = time.Parse("2006-01-02 15:04:05-07:00", *firstActive)
			if err != nil {
				return nil, err
			}
		}

		if redeemedOn != nil {
			gv.RedeemedOn, err = time.Parse("2006-01-02 15:04:05-07:00", *redeemedOn)
			if err != nil {
				return nil, err
			}
		}

		godView = append(godView, gv)
	}

	return godView, nil
}
