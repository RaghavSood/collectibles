package types

import "time"

type Creator struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Slug      string    `json:"slug"`
}
