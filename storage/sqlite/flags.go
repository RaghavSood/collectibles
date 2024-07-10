package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) GetFlags(scope string, key string) ([]types.Flag, error) {
	rows, err := d.db.Query("SELECT flag_scope, flag_type, flag_key FROM flags WHERE flag_scope = ? AND flag_key = ?", scope, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanFlags(rows)
}

func scanFlags(rows *sql.Rows) ([]types.Flag, error) {
	flags := []types.Flag{}
	for rows.Next() {
		var flag types.Flag
		err := rows.Scan(&flag.FlagScope, &flag.FlagType, &flag.FlagKey)
		if err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}
	return flags, nil
}
