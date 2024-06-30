package main

import (
	"flag"
	"os"

	"github.com/RaghavSood/collectibles/storage/sqlite"
	"github.com/RaghavSood/collectibles/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var noindex bool

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})})

	flag.BoolVar(&noindex, "noindex", false, "Don't index the blockchain, run in read-only mode")
	flag.Parse()
}

func main() {
	db, err := sqlite.NewSqliteBackend(noindex)
	defer db.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}

	webServer := web.NewServer(db, noindex)
	webServer.Serve()
}
