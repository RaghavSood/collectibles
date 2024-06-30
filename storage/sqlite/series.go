package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) GetSeries() ([]types.Series, error) {
	rows, err := d.db.Query("SELECT name, created_at, slug FROM series")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSeries(rows)
}

func scanSeries(rows *sql.Rows) ([]types.Series, error) {
	var series []types.Series
	for rows.Next() {
		var s types.Series
		err := rows.Scan(&s.Name, &s.CreatedAt, &s.Slug)
		if err != nil {
			return nil, err
		}
		series = append(series, s)
	}

	return series, nil
}
