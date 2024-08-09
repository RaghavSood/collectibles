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

func (d *SqliteBackend) GetCreatorsBySeries(seriesSlug string) ([]types.Creator, error) {
	rows, err := d.db.Query("SELECT creators.name, creators.created_at, creators.slug FROM creators JOIN series_creators ON creators.slug = series_creators.creator_slug WHERE series_creators.series_slug = ?", seriesSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCreators(rows)
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
