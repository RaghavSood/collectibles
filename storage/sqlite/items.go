package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) GetItems() ([]types.Item, error) {
	rows, err := d.db.Query("SELECT sku, series_slug, serial, created_at FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItems(rows)
}

func (d *SqliteBackend) GetItemPage(pageSize, offset int) ([]types.Item, error) {
	rows, err := d.db.Query("SELECT sku, series_slug, serial, created_at FROM items ORDER BY created_at ASC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItems(rows)
}

func scanItems(rows *sql.Rows) ([]types.Item, error) {
	var items []types.Item
	for rows.Next() {
		var item types.Item
		err := rows.Scan(&item.SKU, &item.SeriesSlug, &item.Serial, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
