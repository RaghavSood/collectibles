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

func (d *SqliteBackend) MarkScriptFastIndex(script, chain string, fastBlockHeight int64) error {
	_, err := d.db.Exec(`UPDATE addresses SET fast_block_height = ? WHERE script = ? AND chain = ?`, fastBlockHeight, script, chain)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`DELETE FROM script_queue WHERE script = ? AND chain = ?`, script, chain)
	return err
}

func (d *SqliteBackend) ScriptExists(script, chain string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM addresses WHERE script = ? AND chain = ?)`, script, chain).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
