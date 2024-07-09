package types

import (
	"encoding/json"
	"time"
)

type ItemAddressSummary struct {
	SKU         string    `json:"sku"`
	Serial      string    `json:"serial"`
	Addresses   string    `json:"addresses"`
	FirstActive time.Time `json:"first_active"`
	RedeemedOn  time.Time `json:"redeemed_on"`
	TotalValue  *BigInt   `json:"total_value"`
	SeriesName  string    `json:"series_name"`
	SeriesSlug  string    `json:"series_slug"`
}

func (is *ItemAddressSummary) SerialString() string {
	if is.Serial == "" {
		return "No Serial"
	}

	return is.Serial
}

func (is *ItemAddressSummary) AddressArray() ([]string, error) {
	result := []string{}

	err := json.Unmarshal([]byte(is.Addresses), &result)

	return result, err
}
