package storage

import "github.com/RaghavSood/collectibles/types"

type Storage interface {
	GetCreators() ([]types.Creator, error)
	GetCreator(slug string) (*types.Creator, error)

	GetSeries() ([]types.Series, error)

	GetItems() ([]types.Item, error)

	GetOnlyScripts(chain string) ([]string, error)
	ScriptExists(script, chain string) (bool, error)

	QueueNewScripts(height int64) error
	GetScriptQueue() ([]types.ScriptQueue, error)
	GetTransactionQueue() ([]types.TransactionQueue, error)
	IncrementScriptQueueTryCount(script, chain string) error
	RecordScriptUnspents(script types.ScriptQueue, unspentTxids []string, unspentHeights []int64) error
	MarkScriptFastIndex(script, chain string, fastBlockHeight int64) error

	RecordTransactionEffects(outpoints []types.Outpoint, spentTxids []string, spentVins []int, spendingTxids []string, spendingVins []int, blockHeight int64, blockTime int) error

	CreatorSummaries() ([]types.CreatorSummary, error)

	SeriesSummaries() ([]types.SeriesSummary, error)
	SeriesSummariesByCreator(creatorSlug string) ([]types.SeriesSummary, error)
	SeriesSummary(seriesSlug string) (*types.SeriesSummary, error)

	ItemSummaries() ([]types.ItemSummary, error)
	ItemSummariesBySeries(seriesSlug string) ([]types.ItemSummary, error)
	ItemSummary(sku string) (*types.ItemSummary, error)

	TransactionSummariesByItem(sku string) ([]types.Transaction, error)
	TransactionSummariesBySeries(seriesSlug string) ([]types.Transaction, error)
}
