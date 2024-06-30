package storage

import "github.com/RaghavSood/collectibles/types"

type Storage interface {
	GetCreators() ([]types.Creator, error)
	GetCreator(slug string) (*types.Creator, error)

	GetSeries() ([]types.Series, error)

	GetItems() ([]types.Item, error)

	QueueNewScripts(height int64) error
	GetScriptQueue() ([]types.ScriptQueue, error)
	IncrementScriptQueueTryCount(script, chain string) error
	RecordScriptUnspents(script types.ScriptQueue, unspentTxids []string, unspentHeights []int64) error
}
