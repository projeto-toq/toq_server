package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) DeleteListingGuarantees(ctx context.Context, tx *sql.Tx, listingVersionID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM guarantees WHERE listing_version_id = ?`

	result, execErr := la.ExecContext(ctx, tx, "delete", query, listingVersionID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.delete_guarantees.exec_error", "error", execErr, "listing_version_id", listingVersionID)
		return fmt.Errorf("exec delete listing guarantees: %w", execErr)
	}

	qty, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.listing.delete_guarantees.rows_affected_error", "error", rowsErr, "listing_version_id", listingVersionID)
		return fmt.Errorf("rows affected for delete listing guarantees: %w", rowsErr)
	}

	if qty == 0 {
		err = fmt.Errorf("no guarantees rows deleted for listing: %w", sql.ErrNoRows)
		// utils.SetSpanError(ctx, err)
		logger.Debug("mysql.listing.delete_guarantees.no_rows", "error", err, "listing_version_id", listingVersionID)
		// return err
	}

	return nil
}
