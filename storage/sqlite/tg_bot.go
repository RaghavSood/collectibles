package sqlite

func (d *SqliteBackend) InsertMessage(chatID int64, message string) error {
	_, err := d.db.Exec("INSERT INTO message_logs (chat_id, message) VALUES (?, ?)", chatID, message)
	return err
}

func (d *SqliteBackend) UpsertTelegramSubscription(chatID int64, scope string, slug string) error {
	_, err := d.db.Exec("INSERT INTO telegram_subscription (chat_id, scope, slug) VALUES (?, ?, ?) ON CONFLICT(chat_id, scope, slug) DO NOTHING", chatID, scope, slug)
	return err
}
