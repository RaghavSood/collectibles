package types

import "time"

type Transaction struct {
	Txid            string    `json:"txid"`
	Vout            int       `json:"vout"`
	Vin             int       `json:"vin"`
	OriginalTxid    string    `json:"original_txid"`
	Value           *BigInt   `json:"value"`
	BlockHeight     int       `json:"block_height"`
	BlockTime       time.Time `json:"block_time"`
	TransactionType string    `json:"transaction_type"`
	SKU             string    `json:"sku"`
	SeriesSlug      string    `json:"series_slug"`
}
