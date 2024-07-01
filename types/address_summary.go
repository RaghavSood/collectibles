package types

import "time"

type AddressSummary struct {
	Address     string    `json:"address"`
	SKU         string    `json:"sku"`
	SeriesSlug  string    `json:"series_slug"`
	Serial      string    `json:"serial"`
	Unspent     int       `json:"unspent"`
	Spent       int       `json:"spent"`
	FirstActive time.Time `json:"first_active"`
	RedeemedOn  time.Time `json:"redeemed_on"`
	TotalValue  *BigInt   `json:"total_value"`
}
