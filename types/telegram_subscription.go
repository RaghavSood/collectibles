package types

type TelegramSubscription struct {
	ChatID int64  `json:"chat_id"`
	Scope  string `json:"scope"`
	Slug   string `json:"slug"`
	Active bool   `json:"active"`
}
