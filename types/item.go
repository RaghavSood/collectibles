package types

import "time"

type Item struct {
	SKU        string    `json:"sku"`
	SeriesSlug string    `json:"series_slug"`
	Serial     string    `json:"serial"`
	CreatedAt  time.Time `json:"created_at"`
}
