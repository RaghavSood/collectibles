package types

import (
	"fmt"
	"time"
)

type Transaction struct {
	Txid            string    `json:"txid"`
	Vout            *int      `json:"vout"`
	Vin             *int      `json:"vin"`
	OriginalTxid    *string   `json:"original_txid"`
	Value           *BigInt   `json:"value"`
	BlockHeight     int       `json:"block_height"`
	BlockTime       time.Time `json:"block_time"`
	TransactionType string    `json:"transaction_type"`
	SKU             string    `json:"sku"`
	SeriesSlug      string    `json:"series_slug"`
	Serial          string    `json:"serial"`
	SeriesName      string    `json:"series_name"`
}

func (t *Transaction) Outpoint() string {
	if t.TransactionType == "outgoing" {
		return fmt.Sprintf("%s:%d", *t.OriginalTxid, *t.Vin)
	}

	return fmt.Sprintf("%s:%d", t.Txid, *t.Vout)
}
