package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) TransactionSummariesByItem(sku string) ([]types.Transaction, error) {
	query := `SELECT txid, vout, vin, original_txid, value, block_height, block_time, transaction_type, sku, series_slug FROM transaction_view WHERE sku = ? ORDER BY block_height DESC;`

	rows, err := d.db.Query(query, sku)
	if err != nil {
		return nil, err
	}

	return scanTransactionSummaries(rows)
}

func (d *SqliteBackend) TransactionSummariesBySeries(slug string) ([]types.Transaction, error) {
	query := `SELECT txid, vout, vin, original_txid, value, block_height, block_time, transaction_type, sku, series_slug FROM transaction_view WHERE series_slug = ? ORDER BY block_height DESC;`

	rows, err := d.db.Query(query, slug)
	if err != nil {
		return nil, err
	}

	return scanTransactionSummaries(rows)
}

func scanTransactionSummaries(rows *sql.Rows) ([]types.Transaction, error) {
	var summaries []types.Transaction
	for rows.Next() {
		var summary types.Transaction
		err := rows.Scan(&summary.Txid, &summary.Vout, &summary.Vin, &summary.OriginalTxid, &summary.Value, &summary.BlockHeight, &summary.BlockTime, &summary.TransactionType, &summary.SKU, &summary.SeriesSlug)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
