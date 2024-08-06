package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) InsertMessage(chatID int64, message string) error {
	_, err := d.db.Exec("INSERT INTO message_logs (chat_id, message) VALUES (?, ?)", chatID, message)
	return err
}

func (d *SqliteBackend) UpsertTelegramSubscription(chatID int64, scope string, slug string) error {
	_, err := d.db.Exec("INSERT INTO telegram_subscription (chat_id, scope, slug) VALUES (?, ?, ?) ON CONFLICT(chat_id, scope, slug) DO NOTHING", chatID, scope, slug)
	return err
}

func (d *SqliteBackend) UnsubscribeTelegram(chatID int64, scope string, slug string) error {
	_, err := d.db.Exec("DELETE FROM telegram_subscription WHERE chat_id = ? AND scope = ? AND slug = ?", chatID, scope, slug)
	return err
}

func (d *SqliteBackend) MatchingTelegramSubscriptions(sku string) ([]types.TelegramSubscription, error) {
	query := `SELECT
  ts.chat_id,
  ts.scope,
  ts.slug,
  ts.active
FROM
  telegram_subscription ts
JOIN
  god_view gv
ON
  (ts.scope = 'item' AND ts.slug = gv.item_id)
  OR (ts.scope = 'series' AND ts.slug = gv.series_id)
  OR (ts.scope = 'creator' AND EXISTS (
    SELECT 1
    FROM json_each(gv.creators) je
    JOIN creators c ON je.value = c.name
    WHERE c.slug = ts.slug
  ))
WHERE
  ts.active = TRUE
  AND gv.item_id = ?;`

	rows, err := d.db.Query(query, sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTelegramSubscriptions(rows)
}

func scanTelegramSubscriptions(rows *sql.Rows) ([]types.TelegramSubscription, error) {
	var subs []types.TelegramSubscription
	for rows.Next() {
		var sub types.TelegramSubscription
		err := rows.Scan(&sub.ChatID, &sub.Scope, &sub.Slug, &sub.Active)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}
