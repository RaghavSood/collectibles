package types

type GeneralStatistics struct {
	Creators      int     `json:"creators"`
	Series        int     `json:"series"`
	Items         int     `json:"items"`
	Addresses     int     `json:"addresses"`
	TotalValue    *BigInt `json:"total_value"`
	TotalRedeemed *BigInt `json:"total_redeemed"`
}
