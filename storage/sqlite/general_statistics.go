package sqlite

import "github.com/RaghavSood/collectibles/types"

func (d *SqliteBackend) GeneralStatistics() (*types.GeneralStatistics, error) {
	query := `SELECT creators, series, items, addresses, total_value, total_redeemed FROM general_statistics;`

	row := d.db.QueryRow(query)

	var stats types.GeneralStatistics
	err := row.Scan(&stats.Creators, &stats.Series, &stats.Items, &stats.Addresses, &stats.TotalValue, &stats.TotalRedeemed)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
