package types

type GradingSlab struct {
	SKU        string `json:"sku"`
	Service    string `json:"service"`
	Identifier string `json:"identifier"`
	Grade      string `json:"grade"`
	ViewLink   string `json:"view_link"`
}
