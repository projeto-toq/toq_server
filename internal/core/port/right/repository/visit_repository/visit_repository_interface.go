package visitrepository

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// VisitRepositoryInterface exposes persistence operations for listing visits.
type VisitRepositoryInterface interface {
	InsertVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) (int64, error)
	UpdateVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) error
	GetVisitByID(ctx context.Context, tx *sql.Tx, id int64) (listingmodel.VisitInterface, error)
	ListVisits(ctx context.Context, tx *sql.Tx, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error)
}
