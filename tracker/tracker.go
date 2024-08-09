package tracker

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/RaghavSood/collectibles/bitcoinrpc"
	btypes "github.com/RaghavSood/collectibles/bitcoinrpc/types"
	"github.com/RaghavSood/collectibles/bloomfilter"
	"github.com/RaghavSood/collectibles/clogger"
	"github.com/RaghavSood/collectibles/electrum"
	"github.com/RaghavSood/collectibles/storage"
	"github.com/RaghavSood/collectibles/tgbot"
	"github.com/RaghavSood/collectibles/types"
	"github.com/RaghavSood/collectibles/util"
)

var log = clogger.NewLogger("tracker")

type Tracker struct {
	db      storage.Storage
	client  *bitcoinrpc.RpcClient
	bf      *bloomfilter.BloomFilter
	eclient *electrum.Electrum
	telebot *tgbot.TgBot
}

func NewTracker(db storage.Storage) *Tracker {
	eclient, err := electrum.NewElectrum()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to electrum server")
	}

	telebot, err := tgbot.NewBot(os.Getenv("TG_BOT_TOKEN"), os.Getenv("TG_CHANNEL_USERNAME"), db)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to telegram bot")
	}

	return &Tracker{
		db:      db,
		client:  bitcoinrpc.NewRpcClient(os.Getenv("BITCOIND_HOST"), os.Getenv("BITCOIND_USER"), os.Getenv("BITCOIND_PASS")),
		bf:      bloomfilter.NewBloomFilter(),
		eclient: eclient,
		telebot: telebot,
	}
}

func (t *Tracker) Run() {
	log.Info().Msg("Starting tracker")

	log.Info().Msg("Loading scripts from database")
	scripts, err := t.db.GetOnlyScripts("bitcoin")
	if err == sql.ErrNoRows {
		log.Info().Msg("No scripts found in database")
		err = nil
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to load scripts")
	} else {
		log.Info().Int("count", len(scripts)).Msg("Loading scripts into bloom filter")
		t.bf.AddStrings(scripts)
		log.Info().Int("count", len(scripts)).Msg("Scripts loaded")
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	dbDumpTicker := time.NewTicker(15 * time.Minute)
	defer dbDumpTicker.Stop()

	err = t.syncTables()
	if err != nil {
		log.Error().Err(err).Msg("Failed to sync tables")
	}

	for {
		select {
		case <-dbDumpTicker.C:
			log.Info().Msg("Updating god view")
			path, err := t.db.UpdateGodView()
			if err != nil {
				log.Error().Err(err).Msg("Failed to update god view")
				continue
			}

			log.Info().
				Str("path", path).
				Msg("God view updated")

			err = t.uploadGodDB(path)
			if err != nil {
				log.Error().Err(err).Msg("Failed to upload god view")
				continue
			}

			// Delete after upload
			err = os.Remove(path)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete god view")
			}

		case <-ticker.C:
			log.Info().Msg("Checking for changes")
			info, err := t.client.GetBlockchainInfo()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get blockchain info")
				continue
			}

			err = t.db.QueueNewScripts(info.Blocks)
			if err != nil {
				log.Error().Err(err).Msg("Failed to queue new scripts")
			}

			t.processScriptQueue()
			t.processTransactionQueue()

			lastBlock, err := t.db.KvGetBlockHeight()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get last block height")
				continue
			}

			// We limit ourselves to batch processing 50 blocks at a time
			// so that other indexing jobs also run often enough
			target := min(lastBlock+1+10, info.Blocks)
			for i := lastBlock + 1; i <= target; i++ {
				blockTime, err := t.processBlock(i)
				if err != nil {
					log.Error().Err(err).Int64("block_height", i).Msg("Failed to process block")
					break
				}

				err = t.syncTables()
				if err != nil {
					log.Error().Err(err).Msg("Failed to sync tables")
				}

				err = t.db.QueueBlockNotification(i, blockTime, "bitcoin")
				if err != nil {
					log.Warn().Err(err).Int64("block_height", i).Msg("Failed to queue block notification")
				}
			}

			t.processBlockNotificationQueue()
		}
	}
}

func (t *Tracker) syncTables() error {
	start := time.Now()
	log.Info().Msg("Syncing computed tables")
	err := t.db.SyncComputedTables()
	if err != nil {
		return fmt.Errorf("failed to sync computed tables: %w", err)
	}
	log.Info().
		Stringer("duration", time.Since(start)).
		Msg("Computed tables synced")

	return nil
}

func (t *Tracker) processBlockNotificationQueue() {
	queuedBlocks, err := t.db.GetBlockNotificationQueue()
	if err == sql.ErrNoRows {
		log.Info().Msg("No blocks in queue")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to get blocks from queue")
		return
	}

	for _, block := range queuedBlocks {
		log.Info().
			Int64("block_height", block.BlockHeight).
			Str("chain", block.Chain).
			Msg("Processing block notification")
		items, err := t.db.RedemptionsByRedeemedOn(block.BlockTime)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get item address summary")
			return
		}

		for _, item := range items {
			log.Info().
				Str("item", item.ItemID).
				Msg("Sending notification")

			start := "An item"
			if item.Serial != nil && *item.Serial != "" {
				start = fmt.Sprintf("Item `%s`", tgbot.EscapeText(*item.Serial))
			}

			series, err := t.db.SeriesSummary(item.SeriesID)
			if err != nil {
				log.Error().Err(err).Str("series", item.SeriesID).Msg("Failed to get series summary")
				return
			}

			message := fmt.Sprintf("%s from series `%s` has been redeemed on %s UTC, worth `%s BTC` \\(%s USD\\)\\.\n\nFirst funded on %s UTC, this item held it's value for *%s*\\.\n\nThere are %d unfunded, %d funded, and %d redeemed items in this series now, worth `%s BTC` \\(%s USD\\)\\.\n\n[View details](https://collectible.money/item/%s)",
				start,
				tgbot.EscapeText(item.SeriesName),
				tgbot.EscapeText(util.ShortUTCTime(item.RedeemedOn)),
				tgbot.EscapeText(item.TotalValue.SatoshisToBTC(true)),
				tgbot.EscapeText(util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(item.TotalValue)))),
				tgbot.EscapeText(util.ShortUTCTime(item.FirstActive)),
				tgbot.EscapeText(util.LifespanString(item.FirstActive, item.RedeemedOn)),
				series.Unfunded,
				series.Unredeemed,
				series.Redeemed,
				tgbot.EscapeText(series.TotalValue.SatoshisToBTC(true)),
				tgbot.EscapeText(util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(series.TotalValue)))),
				tgbot.EscapeText(item.ItemID),
			)
			err = t.telebot.SendMessage(message)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send telegram message")
			}

			matchingSubs, err := t.db.MatchingTelegramSubscriptions(item.ItemID)
			if err != nil {
				log.Error().Err(err).Str("item", item.ItemID).Msg("Failed to get matching subscriptions")
			} else {
				sentItemsToChats := make(map[int64]map[string]bool)

				for _, sub := range matchingSubs {
					if sentItemsToChats[sub.ChatID] != nil && sentItemsToChats[sub.ChatID][item.ItemID] {
						log.Info().
							Int64("chat_id", sub.ChatID).
							Str("scope", sub.Scope).
							Str("slug", sub.Slug).
							Msg("Notification already sent to chat")
						continue
					}

					log.Info().
						Int64("chat_id", sub.ChatID).
						Str("scope", sub.Scope).
						Str("slug", sub.Slug).
						Msg("Sending notification")
					err = t.telebot.SendChatMessage(sub.ChatID, message, false)
					if err != nil {
						log.Error().Err(err).Int64("chat_id", sub.ChatID).Msg("Failed to send telegram message")
						continue
					}

					if sentItemsToChats[sub.ChatID] == nil {
						sentItemsToChats[sub.ChatID] = make(map[string]bool)
					}

					sentItemsToChats[sub.ChatID][item.ItemID] = true
				}
			}
		}

		err = t.db.MarkBlockNotificationProcessed(block.BlockHeight, block.Chain)
		if err != nil {
			log.Error().Err(err).Msg("Failed to mark block notification processed")
			return
		}

		log.Info().
			Int64("block_height", block.BlockHeight).
			Str("chain", block.Chain).
			Msg("Block notification processed")
	}
}

func (t *Tracker) processBlock(height int64) (time.Time, error) {
	log.Info().Int64("block_height", height).Msg("Processing block")
	start := time.Now()

	hash, err := t.client.GetBlockHash(height)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get block hash: %w", err)
	}

	block, err := t.client.GetBlock(hash)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get block: %w", err)
	}

	outpoints, spentTxids, spentVouts, spendingTxids, spendingVins := t.scanTransactions(block.Height, block.Time, block.Tx, "bitcoin")
	err = t.db.RecordTransactionEffects(outpoints, spentTxids, spentVouts, spendingTxids, spendingVins, block.Height, block.Time)
	if err != nil {
		log.Error().Err(err).
			Int64("block_height", height).
			Msg("Failed to record transaction effects")
		return time.Time{}, fmt.Errorf("failed to record transaction effects: %w", err)
	}

	err = t.db.KvSetBlockHeight(height)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to set block height: %w", err)
	}

	log.Info().
		Int64("block_height", height).
		Stringer("duration", time.Since(start)).
		Msg("Block processed")

	return time.Unix(int64(block.Time), 0), nil
}

func (t *Tracker) processTransactionQueue() {
	log.Info().Msg("Processing transaction queue")

	txs, err := t.db.GetTransactionQueue()
	if err == sql.ErrNoRows {
		log.Info().Msg("No transactions in queue")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transactions from queue")
		return
	}

	for _, tx := range txs {
		log.Info().
			Str("txid", tx.Txid).
			Int64("block_height", tx.BlockHeight).
			Msg("Processing transaction")
		txDetails, err := t.client.GetTransaction(tx.Txid)
		if err != nil {
			log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to get transaction details")
			return
		}

		if txDetails.Blockhash == "" {
			log.Info().Str("txid", tx.Txid).Msg("Transaction not yet confirmed")
			return
		}

		block, err := t.client.GetBlock(txDetails.Blockhash)
		if err != nil {
			log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to get block details")
			return
		}

		outpoints, spentTxids, spentVouts, spendingTxids, spendingVins := t.scanTransactions(block.Height, block.Time, []btypes.TransactionDetail{txDetails}, "bitcoin")
		log.Info().
			Int("outpoints", len(outpoints)).
			Int("spent_txids", len(spentTxids)).
			Int("spent_vouts", len(spentVouts)).
			Int("spending_txids", len(spendingTxids)).
			Int("spending_vins", len(spendingVins)).
			Str("txid", tx.Txid).
			Msg("Scanned transactions")

		err = t.db.RecordTransactionEffects(outpoints, spentTxids, spentVouts, spendingTxids, spendingVins, block.Height, block.Time)
		if err != nil {
			log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to record transaction effects")
			return
		}
	}
}

func (t *Tracker) processScriptQueue() {
	log.Info().Msg("Processing script queue")

	scripts, err := t.db.GetScriptQueue()
	if err == sql.ErrNoRows {
		log.Info().Msg("No scripts in queue")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to get addresses from queue")
		return
	}

	for _, script := range scripts {
		if script.Script == "" {
			err = t.db.MarkScriptFastIndex(script.Script, script.Chain, 1)
			continue
		}

		if script.TryCount > 5 {
			log.Warn().Str("script", script.Script).Int("try_count", script.TryCount).Msg("Script has been tried too many times, skipping")
			continue
		}

		log.Info().Str("script", script.Script).Int("try_count", script.TryCount).Msg("Processing script")
		unspentTxids, unspentHeights, err := t.eclient.GetScriptHistory(script.Script)
		if err != nil {
			log.Error().Err(err).Str("script", script.Script).Msg("Failed to get unspents")
			err = t.db.IncrementScriptQueueTryCount(script.Script, script.Chain)
			if err != nil {
				log.Error().Err(err).Str("script", script.Script).Msg("Failed to increment try count")
			}
			continue
		}

		if len(unspentTxids) == 0 {
			log.Info().Str("script", script.Script).Msg("No on-chain activity found for script")
			err = t.db.MarkScriptFastIndex(script.Script, script.Chain, 1)
			continue
		}

		err = t.db.RecordScriptUnspents(script, unspentTxids, unspentHeights)
		if err != nil {
			log.Error().Err(err).Str("script", script.Script).Msg("Failed to record unspents")
			err = t.db.IncrementScriptQueueTryCount(script.Script, script.Chain)
			if err != nil {
				log.Error().Err(err).Str("script", script.Script).Msg("Failed to increment try count")
			}
		}
	}
}

func (t *Tracker) scanTransactions(blockHeight int64, blockTime int, txs []btypes.TransactionDetail, chain string) ([]types.Outpoint, []string, []int, []string, []int) {
	var outpoints []types.Outpoint
	var spentTxids []string
	var spentVouts []int
	var spendingTxids []string
	var spendingVins []int

	for _, tx := range txs {
		for i, vin := range tx.Vin {
			if vin.Coinbase != "" {
				continue
			}

			spentScript := vin.Prevout.ScriptPubKey.Hex

			if t.bf.TestString(spentScript) {
				exists, err := t.db.ScriptExists(spentScript, chain)
				if err != nil {
					log.Error().Err(err).Str("script", spentScript).Msg("Failed to check if script exists")
					continue
				}

				log.Info().
					Str("script", spentScript).
					Bool("exists", exists).
					Str("txid", vin.Txid).
					Int("vout", vin.Vout).
					Msg("Checking if script exists")

				if exists {
					spentTxids = append(spentTxids, vin.Txid)
					spentVouts = append(spentVouts, vin.Vout)
					spendingTxids = append(spendingTxids, tx.Txid)
					spendingVins = append(spendingVins, i)
				}
			}
		}

		for _, vout := range tx.Vout {
			script := vout.ScriptPubKey.Hex

			if t.bf.TestString(script) {
				exists, err := t.db.ScriptExists(script, chain)
				if err != nil {
					log.Error().Err(err).Str("script", script).Msg("Failed to check if script exists")
					continue
				}

				log.Info().
					Str("script", script).
					Bool("exists", exists).
					Str("txid", tx.Txid).
					Int("vout", vout.N).
					Msg("Checking if script exists")

				if exists {
					outpoints = append(outpoints, types.Outpoint{
						Txid:        tx.Txid,
						Vout:        vout.N,
						Script:      script,
						Value:       types.FromBTCString(types.BTCString(vout.Value)),
						BlockHeight: blockHeight,
						BlockTime:   time.Unix(int64(blockTime), 0),
						Chain:       chain,
					})
				}
			}
		}
	}

	return outpoints, spentTxids, spentVouts, spendingTxids, spendingVins
}
