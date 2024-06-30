package types

import "time"

type ScriptQueue struct {
	Script    string    `json:"script"`
	Chain     string    `json:"chain"`
	TryCount  int       `json:"try_count"`
	CreatedAt time.Time `json:"created_at"`
}
