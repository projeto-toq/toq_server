package globalrepository

import (
	"context"
	"database/sql"
)

type GlobalRepoPortInterface interface {
	GetConfiguration(ctx context.Context, tx *sql.Tx) (configuration map[string]string, err error)

	// Transaction related methods
	// StartReadOnlyTransaction starts a database transaction with read-only semantics.
	// It should be used for pure read flows to reduce locking and overhead.
	StartReadOnlyTransaction(ctx context.Context) (tx *sql.Tx, err error)
	StartTransaction(ctx context.Context) (tx *sql.Tx, err error)
	RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) (err error)
}
