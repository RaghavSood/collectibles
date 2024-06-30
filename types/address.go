package types

type Address struct {
	SKU             string `json:"sku"`
	Script          string `json:"script"`
	Address         string `json:"address"`
	Chain           string `json:"chain"`
	FastBlockHeight int64  `json:"fast_block_height"`
}
