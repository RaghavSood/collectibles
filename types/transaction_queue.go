package types

import "time"

type TransactionQueue struct {
	Txid        string    `json:"txid"`
	Chain       string    `json:"chain"`
	BlockHeight int64     `json:"block_height"`
	TryCount    int       `json:"try_count"`
	CreatedAt   time.Time `json:"created_at"`
}
