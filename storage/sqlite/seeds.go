package sqlite

import (
	"encoding/csv"
	"fmt"
)

func (d *SqliteBackend) seedCreators() error {
	csvFile, err := embeddedMigrations.Open("migrations/seeds/creators.csv")
	if err != nil {
		return fmt.Errorf("failed to open creators.csv: %w", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read creators.csv: %w", err)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	_, err = tx.Exec(`CREATE TEMPORARY TABLE creators_temp (name TEXT, slug TEXT)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create temporary table: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO creators_temp (name, slug) VALUES (?, ?)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, record := range records[1:] {
		_, err := stmt.Exec(record[0], record[1])
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert creator: %w", err)
		}
	}

	_, err = tx.Exec(`INSERT INTO creators (name, slug) SELECT name, slug FROM creators_temp WHERE slug NOT IN (SELECT slug FROM creators)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert creators: %w", err)
	}

	return tx.Commit()
}
