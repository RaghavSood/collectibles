package sqlite

import (
	"database/sql"
	"fmt"
)

func (d *SqliteBackend) SyncComputedTables() error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	err = d.syncTransactionSummary(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to sync transaction summary: %w", err)
	}

	err = d.syncItemSummary(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to sync item summary: %w", err)
	}

	err = d.syncAddressSummary(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to sync address summary: %w", err)
	}

	return tx.Commit()
}

func (d *SqliteBackend) syncTransactionSummary(tx *sql.Tx) error {
	// Truncate existing transaction_summary_c table
	_, err := tx.Exec(`DELETE FROM transaction_summary_c`)
	if err != nil {
		return fmt.Errorf("failed to truncate transaction_summary_c: %w", err)
	}

	// Insert new data into transaction_summary_c
	_, err = tx.Exec(`INSERT INTO transaction_summary_c SELECT * FROM transaction_view`)
	if err != nil {
		return fmt.Errorf("failed to insert into transaction_summary_c: %w", err)
	}

	return nil
}

func (d *SqliteBackend) syncItemSummary(tx *sql.Tx) error {
	// Truncate existing item_summary_c table
	_, err := tx.Exec(`DELETE FROM item_summary_c`)
	if err != nil {
		return fmt.Errorf("failed to truncate item_summary_c: %w", err)
	}

	// Insert new data into item_summary_c
	_, err = tx.Exec(`INSERT INTO item_summary_c SELECT * FROM item_summary`)
	if err != nil {
		return fmt.Errorf("failed to insert into item_summary_c: %w", err)
	}

	return nil
}

func (d *SqliteBackend) syncAddressSummary(tx *sql.Tx) error {
	// Truncate existing address_summary_c table
	_, err := tx.Exec(`DELETE FROM address_summary_c`)
	if err != nil {
		return fmt.Errorf("failed to truncate address_summary_c: %w", err)
	}

	// Insert new data into address_summary_c
	_, err = tx.Exec(`INSERT INTO address_summary_c SELECT * FROM address_summary`)
	if err != nil {
		return fmt.Errorf("failed to insert into address_summary_c: %w", err)
	}

	return nil
}
