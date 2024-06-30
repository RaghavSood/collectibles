package tracker

import (
	"database/sql"
	"os"
	"time"

	"github.com/RaghavSood/collectibles/bitcoinrpc"
	"github.com/RaghavSood/collectibles/clogger"
	"github.com/RaghavSood/collectibles/electrum"
	"github.com/RaghavSood/collectibles/storage"
)

var log = clogger.NewLogger("tracker")

type Tracker struct {
	db      storage.Storage
	client  *bitcoinrpc.RpcClient
	eclient *electrum.Electrum
}

func NewTracker(db storage.Storage) *Tracker {
	eclient, err := electrum.NewElectrum()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to electrum server")
	}

	return &Tracker{
		db:      db,
		client:  bitcoinrpc.NewRpcClient(os.Getenv("BITCOIND_HOST"), os.Getenv("BITCOIND_USER"), os.Getenv("BITCOIND_PASS")),
		eclient: eclient,
	}
}

func (t *Tracker) Run() {
	log.Info().Msg("Starting tracker")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
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
