package sqlite

func (d *SqliteBackend) InsertMessage(chatID int64, message string) error {
	_, err := d.db.Exec("INSERT INTO message_logs (chat_id, message) VALUES (?, ?)", chatID, message)
	return err
}
