package mysqllistingadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (la *ListingAdapter) DeleteListingFinancingBlockers(ctx context.Context, tx *sql.Tx, listingID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `DELETE FROM financing_blockers WHERE listing_id = ?`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteListingListingBlockers: error preparing statement", "error", err)
		err = status.Error(codes.Internal, "Internal server error")
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, listingID)
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteListingListingBlockers: error executing statement", "error", err)
		err = status.Error(codes.Internal, "Internal server error")
		return
	}

	qty, err := result.RowsAffected()
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteListingListingBlockers: error getting rows affected", "error", err)
		err = status.Error(codes.Internal, "Internal server error")
		return
	}

	if qty == 0 {
		err = status.Error(codes.NotFound, "Listing blocker not found")
		return
	}

	return
}
