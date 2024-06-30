package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) QueueNewScripts(_ int64) error {
	_, err := d.db.Exec(`INSERT INTO script_queue (script, chain) SELECT script, chain FROM addresses WHERE fast_block_height = 0 ON CONFLICT DO NOTHING`)
	return err
}

func (d *SqliteBackend) GetScriptQueue() ([]types.ScriptQueue, error) {
	rows, err := d.db.Query(`SELECT script, chain, try_count, created_at FROM script_queue`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanScriptQueue(rows)
}

func (d *SqliteBackend) GetTransactionQueue() ([]types.TransactionQueue, error) {
	rows, err := d.db.Query(`SELECT txid, chain, block_height, try_count, created_at FROM transaction_queue ORDER BY block_height ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTransactionQueue(rows)
}

func (d *SqliteBackend) IncrementScriptQueueTryCount(script, chain string) error {
	_, err := d.db.Exec(`UPDATE script_queue SET try_count = try_count + 1 WHERE script = ? AND chain = ?`, script, chain)
	return err
}

func (d *SqliteBackend) RecordScriptUnspents(script types.ScriptQueue, unspentTxids []string, unspentHeights []int64) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	for i, txid := range unspentTxids {
		_, err = tx.Exec(`INSERT INTO transaction_queue (txid, chain, block_height) VALUES (?, ?, ?) ON CONFLICT DO NOTHING`, txid, script.Chain, unspentHeights[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(`DELETE FROM script_queue WHERE script = ? AND chain = ?`, script.Script, script.Chain)
	if err != nil {
		tx.Rollback()
		return err
	}

	maxHeight := int64(0)
	for _, height := range unspentHeights {
		if height > maxHeight {
			maxHeight = height
		}
	}

	_, err = tx.Exec(`UPDATE addresses SET fast_block_height = ? WHERE script = ? AND chain = ?`, maxHeight, script.Script, script.Chain)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func scanScriptQueue(rows *sql.Rows) ([]types.ScriptQueue, error) {
	var scripts []types.ScriptQueue
	for rows.Next() {
		var script types.ScriptQueue
		err := rows.Scan(&script.Script, &script.Chain, &script.TryCount, &script.CreatedAt)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}
	return scripts, nil
}

func scanTransactionQueue(rows *sql.Rows) ([]types.TransactionQueue, error) {
	var transactions []types.TransactionQueue
	for rows.Next() {
		var tx types.TransactionQueue
		err := rows.Scan(&tx.Txid, &tx.Chain, &tx.BlockHeight, &tx.TryCount, &tx.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
