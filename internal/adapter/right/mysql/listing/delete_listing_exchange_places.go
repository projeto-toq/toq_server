package mysqllistingadapter

import (
	"context"
	"database/sql"
	"log/slog"

	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) DeleteListingExchangePlaces(ctx context.Context, tx *sql.Tx, listingID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `DELETE FROM exchange_places WHERE listing_id = ?`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteExchangePlaces: error preparing statement", "error", err)
		err = utils.ErrInternalServer
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, listingID)
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteExchangePlaces: error executing statement", "error", err)
		err = utils.ErrInternalServer
		return
	}

	qty, err := result.RowsAffected()
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteExchangePlaces: error getting rows affected", "error", err)
		err = utils.ErrInternalServer
		return
	}

	if qty == 0 {
		err = utils.ErrInternalServer
		return
	}

	return
}
