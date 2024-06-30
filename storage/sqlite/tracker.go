package sqlite

import "github.com/RaghavSood/collectibles/types"

func (d *SqliteBackend) RecordTransactionEffects(outpoints []types.Outpoint, spentTxids []string, spentVins []int, spendingTxids []string, spendingVins []int, blockHeight int64, blockTime int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	for _, outpoint := range outpoints {
		_, err = tx.Exec(`INSERT INTO outpoints (txid, vout, script, value, block_height, block_time, chain) VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT (txid, vout) DO NOTHING`, outpoint.Txid, outpoint.Vout, outpoint.Script, outpoint.Value.String(), outpoint.BlockHeight, outpoint.BlockTime, outpoint.Chain)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for i, txid := range spentTxids {
		_, err := tx.Exec(`UPDATE outpoints SET spending_txid = ?, spending_vin = ?, spending_block_height = ?, spending_block_time = ? WHERE txid = ? AND vout = ?`, spendingTxids[i], spendingVins[i], blockHeight, blockTime, txid, spentVins[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, outpoint := range outpoints {
		_, err := tx.Exec(`DELETE FROM transaction_queue WHERE txid = ?`, outpoint.Txid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, txid := range spendingTxids {
		_, err := tx.Exec(`DELETE FROM transaction_queue WHERE txid = ?`, txid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
