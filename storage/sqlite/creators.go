package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) GetCreators() ([]types.Creator, error) {
	rows, err := d.db.Query("SELECT name, created_at, slug FROM creators")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCreators(rows)
}

func (d *SqliteBackend) GetCreator(slug string) (*types.Creator, error) {
	row := d.db.QueryRow("SELECT name, created_at, slug FROM creators WHERE slug = ?", slug)

	var creator types.Creator
	err := row.Scan(&creator.Name, &creator.CreatedAt, &creator.Slug)
	if err != nil {
		return nil, err
	}

	return &creator, nil
}

func scanCreators(rows *sql.Rows) ([]types.Creator, error) {
	var creators []types.Creator
	for rows.Next() {
		var creator types.Creator
		err := rows.Scan(&creator.Name, &creator.CreatedAt, &creator.Slug)
		if err != nil {
			return nil, err
		}
		creators = append(creators, creator)
	}

	return creators, nil
}
