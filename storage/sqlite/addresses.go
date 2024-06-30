package sqlite

func (d *SqliteBackend) GetOnlyScripts(chain string) ([]string, error) {
	rows, err := d.db.Query(`SELECT script FROM addresses WHERE chain = ?`, chain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []string
	for rows.Next() {
		var script string
		err = rows.Scan(&script)
		if err != nil {
			return nil, err
		}

		scripts = append(scripts, script)
	}

	return scripts, nil
}

func (d *SqliteBackend) ScriptExists(script, chain string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM addresses WHERE script = ? AND chain = ?)`, script, chain).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
