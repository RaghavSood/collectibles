package types

import "time"

type BlockNotificationQueue struct {
	BlockHeight int64     `json:"block_height"`
	Chain       string    `json:"chain"`
	CreatedAt   time.Time `json:"created_at"`
}
