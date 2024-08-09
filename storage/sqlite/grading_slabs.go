package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) GetGradingSlabsByItem(sku string) ([]types.GradingSlab, error) {
	rows, err := d.db.Query("SELECT sku, service, identifier, grade, view_link FROM grading_slabs WHERE sku = ?", sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanGradingSlabs(rows)
}

func scanGradingSlabs(rows *sql.Rows) ([]types.GradingSlab, error) {
	var gradingSlabs []types.GradingSlab
	for rows.Next() {
		var gradingSlab types.GradingSlab
		err := rows.Scan(&gradingSlab.SKU, &gradingSlab.Service, &gradingSlab.Identifier, &gradingSlab.Grade, &gradingSlab.ViewLink)
		if err != nil {
			return nil, err
		}
		gradingSlabs = append(gradingSlabs, gradingSlab)
	}
	return gradingSlabs, nil
}
