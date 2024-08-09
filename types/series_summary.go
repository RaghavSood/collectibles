package types

type SeriesSummary struct {
	Slug       string    `json:"slug"`
	Name       string    `json:"name"`
	ItemCount  int       `json:"item_count"`
	TotalValue *BigInt   `json:"total_value"`
	Unfunded   int       `json:"unfunded"`
	Unredeemed int       `json:"unredeemed"`
	Redeemed   int       `json:"redeemed"`
	Creators   []Creator `json:"creators"`
}
