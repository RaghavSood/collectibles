package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) GodView() ([]types.GodView, error) {
	rows, err := d.db.Query(`
		SELECT series_name, series_id, creators, item_id, serial, addresses, total_value, first_active, redeemed_on, balance
		FROM god_view
		ORDER BY series_id, item_id, serial;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query god view: %w", err)
	}
	defer rows.Close()

	return scanGodView(rows)
}

func (d *SqliteBackend) Search(query string) ([]types.GodView, error) {
	rows, err := d.db.Query(`
    WITH slab_matches AS (
		  SELECT sku FROM grading_slabs WHERE identifier LIKE '%' || $1 || '%'
		)
		SELECT series_name, series_id, creators, item_id, serial, addresses, total_value, first_active, redeemed_on, balance
		FROM god_view
		WHERE series_name LIKE '%' || $1 || '%' OR addresses LIKE '%' || $1 || '%' OR serial LIKE '%' || $1 || '%' OR creators LIKE '%' || $1 || '%' OR item_id IN slab_matches
		ORDER BY series_id, item_id, serial;
	`, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query god view: %w", err)
	}
	defer rows.Close()

	return scanGodView(rows)
}

func (d *SqliteBackend) RecentRedemptions(limit int) ([]types.GodView, error) {
	rows, err := d.db.Query(`
		SELECT series_name, series_id, creators, item_id, serial, addresses, total_value, first_active, redeemed_on, balance
		FROM god_view
		WHERE redeemed_on IS NOT NULL
		ORDER BY redeemed_on DESC
		LIMIT $1;
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query god view: %w", err)
	}
	defer rows.Close()

	return scanGodView(rows)
}

func (d *SqliteBackend) StolenLostItems() ([]types.GodView, error) {
	rows, err := d.db.Query(`
		SELECT series_name, series_id, creators, item_id, serial, addresses, total_value, first_active, redeemed_on, balance
		FROM god_view
		WHERE item_id IN (SELECT flag_key FROM flags WHERE flag_scope = 'items' AND flag_type = 'stolen')
		ORDER BY series_id, item_id, serial;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query god view: %w", err)
	}
	defer rows.Close()

	return scanGodView(rows)
}

func (d *SqliteBackend) RedemptionsByRedeemedOn(redeemedOn time.Time) ([]types.GodView, error) {
	rows, err := d.db.Query(`
		SELECT series_name, series_id, creators, item_id, serial, addresses, total_value, first_active, redeemed_on, balance
		FROM god_view
		WHERE redeemed_on = $1
		ORDER BY series_id, item_id, serial;
	`, redeemedOn)
	if err != nil {
		return nil, fmt.Errorf("failed to query god view: %w", err)
	}
	defer rows.Close()

	return scanGodView(rows)
}

func scanGodView(rows *sql.Rows) ([]types.GodView, error) {
	var godView []types.GodView
	for rows.Next() {
		var gv types.GodView
		var firstActive *string
		var redeemedOn *string
		err := rows.Scan(&gv.SeriesName, &gv.SeriesID, &gv.Creators, &gv.ItemID, &gv.Serial, &gv.Addresses, &gv.TotalValue, &firstActive, &redeemedOn, &gv.Balance)
		if err != nil {
			return nil, fmt.Errorf("failed to scan god view: %w", err)
		}

		if firstActive != nil {
			gv.FirstActive, err = parseDbTimeString(firstActive)
			if err != nil {
				return nil, err
			}
		}

		if redeemedOn != nil {
			gv.RedeemedOn, err = parseDbTimeString(redeemedOn)
			if err != nil {
				return nil, err
			}
		}

		godView = append(godView, gv)
	}

	return godView, nil
}
