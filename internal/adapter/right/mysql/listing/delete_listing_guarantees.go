package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) DeleteListingGuarantees(ctx context.Context, tx *sql.Tx, listingID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM guarantees WHERE listing_id = ?`
	defer la.ObserveOnComplete("delete", query)()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.delete_guarantees.prepare_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("prepare delete listing guarantees: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, listingID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.delete_guarantees.exec_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("exec delete listing guarantees: %w", err)
	}

	qty, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.delete_guarantees.rows_affected_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("rows affected for delete listing guarantees: %w", err)
	}

	if qty == 0 {
		err = fmt.Errorf("no guarantees rows deleted for listing: %w", sql.ErrNoRows)
		// utils.SetSpanError(ctx, err)
		logger.Debug("mysql.listing.delete_guarantees.no_rows", "error", err, "listing_id", listingID)
		// return err
	}

	return nil
}
