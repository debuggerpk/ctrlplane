package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	Config = config.Config
)

var (
	DefaultConfig = config.Default
)

func WithConfig(conf *Config) config.ConfigOption {
	return config.WithConfig(conf)
}

// Get is a wrapper around the dbcfg.Instance singleton.
func Get(opts ...config.ConfigOption) *config.Config {
	return config.Instance(opts...)
}

// Queries is a wrapper around the dbcfg.Queries singleton.
func Queries() *entities.Queries {
	return config.Queries()
}

// Transaction begins the transaction and wraps the queries in a transaction.
//
// Example:
//
//	tx, qtx, err := db.Transaction(ctx)
//	if err != nil { return err }
//
//	defer func() { _ = tx.Rollback(ctx) }()
//
//	// Do something with qtx. Any time you return on error, the transaction will be rolled back.
//	...
//
//	// Commit the transaction.
//	err = tx.Commit(ctx)
//	if err != nil { return err }
//
//	return nil
func Transaction(ctx context.Context) (pgx.Tx, *entities.Queries, error) {
	tx, err := Get().Pool().Begin(ctx)
	if err != nil {
		slog.Error("db: error creating transaction ...", "error", err.Error())

		return nil, nil, err
	}

	return tx, Queries().WithTx(tx), nil
}
