package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) DeleteListingFeatures(ctx context.Context, tx *sql.Tx, listingID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM features WHERE listing_id = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.delete_features.prepare_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("prepare delete listing features: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, listingID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.delete_features.exec_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("exec delete listing features: %w", err)
	}

	qty, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.delete_features.rows_affected_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("rows affected for delete listing features: %w", err)
	}

	if qty == 0 {
		err = fmt.Errorf("no features rows deleted for listing: %w", sql.ErrNoRows)
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.delete_features.no_rows", "error", err, "listing_id", listingID)
		return err
	}

	return nil
}
