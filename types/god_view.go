package types

import "time"

type GodView struct {
	SeriesName  string    `json:"series_name"`
	SeriesID    string    `json:"series_id"`
	Creators    string    `json:"creators"`
	ItemID      string    `json:"item_id"`
	Serial      *string   `json:"serial"`
	Address     *string   `json:"address"`
	TotalValue  *BigInt   `json:"total_value"`
	FirstActive time.Time `json:"first_active"`
	RedeemedOn  time.Time `json:"redeemed_on"`
}

func (g GodView) SerialString() string {
	if g.Serial == nil || *g.Serial == "" {
		return "No Serial"
	}
	return *g.Serial
}
