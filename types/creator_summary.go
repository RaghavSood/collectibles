package types

type CreatorSummary struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	SeriesCount int     `json:"series_count"`
	ItemCount   int     `json:"item_count"`
	TotalValue  *BigInt `json:"total_value"`
	Unfunded    int     `json:"unfunded"`
	Redeemed    int     `json:"redeemed"`
	Unredeemed  int     `json:"unredeemed"`
}
