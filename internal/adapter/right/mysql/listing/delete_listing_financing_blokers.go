package mysqllistingadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
		err = fmt.Errorf("prepare delete financing blockers: %w", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, listingID)
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteListingListingBlockers: error executing statement", "error", err)
		err = fmt.Errorf("exec delete financing blockers: %w", err)
		return
	}

	qty, err := result.RowsAffected()
	if err != nil {
		slog.Error("mysqllistingadapter/DeleteListingListingBlockers: error getting rows affected", "error", err)
		err = fmt.Errorf("rows affected for delete financing blockers: %w", err)
		return
	}

	if qty == 0 {
		err = errors.New("no financing_blockers rows deleted for listing")
		return
	}

	return
}
