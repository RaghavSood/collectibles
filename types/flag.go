package types

const (
	FLAG_SCOPE_CREATORS = "creators"
	FLAG_SCOPE_SERIES   = "series"
	FLAG_SCOPE_ITEMS    = "items"
)

type Flag struct {
	FlagScope string `json:"flag_scope"`
	FlagType  string `json:"flag_type"`
	FlagKey   string `json:"flag_key"`
}
