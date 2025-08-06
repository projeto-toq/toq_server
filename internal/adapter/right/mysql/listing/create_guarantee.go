package mysqllistingadapter

import (
	"context"
	"database/sql"
	"log/slog"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (la *ListingAdapter) CreateGuarantee(ctx context.Context, tx *sql.Tx, guarantee listingmodel.GuaranteeInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO guarantees (listing_id, priority, guarantee) VALUES (?, ?, ?);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("mysqllistingadapter/CreateGuarantee: error preparing statement", "error", err)
		err = status.Error(codes.Internal, "Failed to prepare statement")
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, guarantee.ListingID(), guarantee.Priority(), guarantee.Guarantee())
	if err != nil {
		slog.Error("mysqllistingadapter/CreateGuarantee: error executing statement", "error", err)
		err = status.Error(codes.Internal, "Failed to execute statement")
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("mysqllistingadapter/CreateGuarantee: error getting last insert ID", "error", err)
		err = status.Error(codes.Internal, "Failed to retrieve last insert ID")
		return
	}

	guarantee.SetID(id)

	return
}
