package visitrepository

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// VisitRepositoryInterface exposes persistence operations for listing visits.
//
// Contract notes:
//   - All methods participate in optional transactions (tx can be nil for reads when allowed).
//   - Return sql.ErrNoRows when the target visit is not found (services map to domain/HTTP).
//   - Implementations must use InstrumentedAdapter for tracing/metrics and avoid SELECT * per guide.
type VisitRepositoryInterface interface {
	// InsertVisit persists a new visit and returns its auto-generated ID.
	// Must run inside a transaction; returns infrastructure errors or validation errors from DB constraints.
	InsertVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) (int64, error)

	// UpdateVisit updates mutable fields of an existing visit; sql.ErrNoRows when PK not found.
	UpdateVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) error

	// GetVisitByID retrieves a visit by primary key; returns sql.ErrNoRows if absent.
	GetVisitByID(ctx context.Context, tx *sql.Tx, id int64) (listingmodel.VisitInterface, error)

	// ListVisits lists visits with filtering/pagination and hydrates the active listing snapshot for each row.
	ListVisits(ctx context.Context, tx *sql.Tx, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error)
}
