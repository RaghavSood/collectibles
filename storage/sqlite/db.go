package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/RaghavSood/collectibles/clogger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*
var embeddedMigrations embed.FS

var log = clogger.NewLogger("sqlite")

type SqliteBackend struct {
	db *sql.DB
}

func NewSqliteBackend(readonly bool) (*SqliteBackend, error) {
	path := os.Getenv("DB_PATH")
	if path == "" {
		return nil, fmt.Errorf("DB_PATH environment variable must be set")
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better performance
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}
	// WAL prefers synchronous=normal
	_, err = db.Exec("PRAGMA synchronous=NORMAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to set synchronous mode: %w", err)
	}
	// Set a large cache size
	_, err = db.Exec("PRAGMA cache_size=-2048000;")
	if err != nil {
		return nil, fmt.Errorf("failed to set cache size: %w", err)
	}
	// Increase busy_timeout
	_, err = db.Exec("PRAGMA busy_timeout=30000;")
	if err != nil {
		return nil, fmt.Errorf("failed to set busy_timeout: %w", err)
	}

	log.Info().
		Str("path", path).
		Msg("Database opened")

	backend := &SqliteBackend{db: db}
	if !readonly {
		if err := backend.Migrate(); err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}

		err = backend.seedCreators()
		if err != nil {
			return nil, fmt.Errorf("failed to seed creators: %w", err)
		}

		err = backend.seedSeries()
		if err != nil {
			return nil, fmt.Errorf("failed to seed series: %w", err)
		}

		err = backend.seedItems()
		if err != nil {
			return nil, fmt.Errorf("failed to seed items: %w", err)
		}

		err = backend.seedAddresses()
		if err != nil {
			return nil, fmt.Errorf("failed to seed addresses: %w", err)
		}

		err = backend.seedFlags()
		if err != nil {
			return nil, fmt.Errorf("failed to seed flags: %w", err)
		}

		err = backend.seedGradingSlabs()
		if err != nil {
			return nil, fmt.Errorf("failed to seed grading slabs: %w", err)
		}
	}

	return backend, nil
}

func (d *SqliteBackend) Close() error {
	return d.db.Close()
}

func (d *SqliteBackend) Migrate() error {
	goose.SetBaseFS(embeddedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(d.db, "migrations"); err != nil {
		return fmt.Errorf("failed to run goose up: %w", err)
	}
	return nil
}
