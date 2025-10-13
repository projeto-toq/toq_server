package globalrepository

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

type GlobalRepoPortInterface interface {
	CreateAudit(ctx context.Context, tx *sql.Tx, audit globalmodel.AuditInterface) (err error)

	GetConfiguration(ctx context.Context, tx *sql.Tx) (configuration map[string]string, err error)

	ListCatalogValues(ctx context.Context, tx *sql.Tx, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error)
	GetCatalogValueByID(ctx context.Context, tx *sql.Tx, category string, id uint8) (listingmodel.CatalogValueInterface, error)
	GetCatalogValueBySlug(ctx context.Context, tx *sql.Tx, category, slug string) (listingmodel.CatalogValueInterface, error)
	GetNextCatalogValueID(ctx context.Context, tx *sql.Tx, category string) (uint8, error)
	CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	UpdateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	SoftDeleteCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8) error

	// Transaction related methods
	// StartReadOnlyTransaction starts a database transaction with read-only semantics.
	// It should be used for pure read flows to reduce locking and overhead.
	StartReadOnlyTransaction(ctx context.Context) (tx *sql.Tx, err error)
	StartTransaction(ctx context.Context) (tx *sql.Tx, err error)
	RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) (err error)
}
