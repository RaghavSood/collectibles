package types

type ItemSummary struct {
	SKU           string  `json:"sku"`
	Serial        string  `json:"serial"`
	SeriesName    string  `json:"series_name"`
	SeriesSlug    string  `json:"series_slug"`
	TotalValue    *BigInt `json:"total_value"`
	Unspent       int     `json:"unspent"`
	Spent         int     `json:"spent"`
	TotalReceived *BigInt `json:"total_received"`
	TotalSpent    *BigInt `json:"total_spent"`
	Unfunded      int     `json:"unfunded"`
	Unredeemed    int     `json:"unredeemed"`
	Redeemed      int     `json:"redeemed"`
}

func (is *ItemSummary) SerialString() string {
	if is.Serial == "" {
		return "No Serial"
	}

	return is.Serial
}
