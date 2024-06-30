package types

import "time"

type Outpoint struct {
	Txid                string    `json:"txid"`
	Vout                int       `json:"vout"`
	Script              string    `json:"script"`
	Value               *BigInt   `json:"value"`
	BlockHeight         int64     `json:"block_height"`
	BlockTime           time.Time `json:"block_time"`
	SpendingTxid        string    `json:"spending_txid"`
	SpendingVin         int       `json:"spending_vin"`
	SpendingBlockHeight int64     `json:"spending_block_height"`
	SpendingBlockTime   time.Time `json:"spending_block_time"`
}
