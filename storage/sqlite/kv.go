package sqlite

func (s *SqliteBackend) KvGetBlockHeight() (int64, error) {
	var height int64
	err := s.db.QueryRow(`SELECT value FROM kv WHERE key = 'block_height'`).Scan(&height)
	return height, err
}

func (s *SqliteBackend) KvSetBlockHeight(height int64) error {
	_, err := s.db.Exec(`INSERT INTO kv (key, value) VALUES ('block_height', ?) ON CONFLICT(key) DO UPDATE SET value = ?`, height, height)
	return err
}
