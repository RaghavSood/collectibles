package storage

import (
	"time"

	"github.com/RaghavSood/collectibles/types"
)

type Storage interface {
	GetCreators() ([]types.Creator, error)
	GetCreator(slug string) (*types.Creator, error)

	GetSeries() ([]types.Series, error)

	GetItems() ([]types.Item, error)
	GetItemPage(pageSize, offset int) ([]types.Item, error)

	GetOnlyScripts(chain string) ([]string, error)
	ScriptExists(script, chain string) (bool, error)

	QueueNewScripts(height int64) error
	GetScriptQueue() ([]types.ScriptQueue, error)
	GetTransactionQueue() ([]types.TransactionQueue, error)
	IncrementScriptQueueTryCount(script, chain string) error
	RecordScriptUnspents(script types.ScriptQueue, unspentTxids []string, unspentHeights []int64) error
	MarkScriptFastIndex(script, chain string, fastBlockHeight int64) error
	GetQueueStats() (int, int, error)

	RecordTransactionEffects(outpoints []types.Outpoint, spentTxids []string, spentVins []int, spendingTxids []string, spendingVins []int, blockHeight int64, blockTime int) error

	QueueBlockNotification(height int64, blockTime time.Time, chain string) error
	GetBlockNotificationQueue() ([]types.BlockNotificationQueue, error)

	CreatorSummaries() ([]types.CreatorSummary, error)
	CreatorSummary(creatorSlug string) (*types.CreatorSummary, error)

	SeriesSummaries() ([]types.SeriesSummary, error)
	SeriesSummariesByCreator(creatorSlug string) ([]types.SeriesSummary, error)
	SeriesSummary(seriesSlug string) (*types.SeriesSummary, error)

	ItemSummaries() ([]types.ItemSummary, error)
	ItemSummariesBySeries(seriesSlug string) ([]types.ItemSummary, error)
	ItemSummary(sku string) (*types.ItemSummary, error)

	TransactionSummariesByItem(sku string) ([]types.Transaction, error)
	TransactionSummariesBySeries(seriesSlug string, limit int) ([]types.Transaction, error)
	TransactionSummariesByCreator(creatorSlug string, limit int) ([]types.Transaction, error)
	TransactionSummaries(limit int) ([]types.Transaction, error)

	AddressSummariesByItem(sku string) ([]types.AddressSummary, error)
	AddressSummariesBySeries(seriesSlug string) ([]types.AddressSummary, error)

	ItemAddressSummariesBySeries(seriesSlug string) ([]types.ItemAddressSummary, error)
	ItemAddressSummariesByRedeemedOn(redeemedOn time.Time) ([]types.ItemAddressSummary, error)
	MarkBlockNotificationProcessed(height int64, chain string) error

	GeneralStatistics() (*types.GeneralStatistics, error)

	GodView() ([]types.GodView, error)
	Search(query string) ([]types.GodView, error)
	RecentRedemptions(limit int) ([]types.GodView, error)
	RedemptionsByRedeemedOn(redeemedOn time.Time) ([]types.GodView, error)
	UpdateGodView() (string, error)

	KvGetBlockHeight() (int64, error)
	KvSetBlockHeight(height int64) error

	GetFlags(scope string, key string) ([]types.Flag, error)

	SyncComputedTables() error

	InsertMessage(chatID int64, message string) error
	UpsertTelegramSubscription(chatID int64, scope string, slug string) error
	UnsubscribeTelegram(chatID int64, scope string, slug string) error
}
