package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (d *SqliteBackend) UpdateGodView() (string, error) {
	pathPrefix := os.Getenv("GOD_DB_PATH")
	dbFile := fmt.Sprintf("collectibles_%d.sqlite", time.Now().Unix())
	newDbPath := filepath.Join(pathPrefix, dbFile)

	tx, err := d.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}

	_, err = tx.Exec("ATTACH DATABASE ? AS goddb", newDbPath)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to attach database: %w", err)
	}

	_, err = tx.Exec(`CREATE TABLE goddb.god_view(
										  series_name TEXT,
										  series_id TEXT,
										  creators TEXT,
										  item_id TEXT,
										  serial TEXT,
										  addresses TEXT,
										  total_value INTEGER,
										  first_active DATETIME,
										  redeemed_on DATETIME,
											balance INTEGER
										);`)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create table: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO goddb.god_view SELECT * FROM main.god_view;`)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to insert into table: %w", err)
	}

	_, err = tx.Exec(`
		CREATE INDEX goddb.god_view_series_id ON god_view(series_id);
    CREATE INDEX goddb.god_view_creators ON god_view(creators);
		CREATE INDEX goddb.god_view_item_id ON god_view(item_id);
		CREATE INDEX goddb.god_view_serial ON god_view(serial);
		CREATE INDEX goddb.god_view_addresses ON god_view(addresses);
		CREATE INDEX goddb.god_view_total_value ON god_view(total_value);
		CREATE INDEX goddb.god_view_first_active ON god_view(first_active);
		CREATE INDEX goddb.god_view_redeemed_on ON god_view(redeemed_on);
		CREATE INDEX goddb.god_view_balance ON god_view(balance);
	`)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create indexes: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	_, err = d.db.Exec("DETACH DATABASE goddb;")
	if err != nil {
		return "", fmt.Errorf("failed to detach database: %w", err)
	}

	return newDbPath, nil
}
