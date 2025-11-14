package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteListingWarehouseAdditionalFloors removes all warehouse additional floors for a listing version.
// This is typically called before inserting updated floors.
func (la *ListingAdapter) DeleteListingWarehouseAdditionalFloors(ctx context.Context, tx *sql.Tx, listingVersionID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM warehouse_additional_floors WHERE listing_version_id = ?`

	result, execErr := la.ExecContext(ctx, tx, "delete", query, listingVersionID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.delete_warehouse_additional_floors.exec_error", "error", execErr, "listing_version_id", listingVersionID)
		return fmt.Errorf("exec delete listing warehouse additional floors: %w", execErr)
	}

	qty, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.listing.delete_warehouse_additional_floors.rows_affected_error", "error", rowsErr, "listing_version_id", listingVersionID)
		return fmt.Errorf("rows affected for delete listing warehouse additional floors: %w", rowsErr)
	}

	if qty == 0 {
		err = fmt.Errorf("no warehouse additional floors rows deleted for listing: %w", sql.ErrNoRows)
		logger.Debug("mysql.listing.delete_warehouse_additional_floors.no_rows", "error", err, "listing_version_id", listingVersionID)
	}

	return nil
}
