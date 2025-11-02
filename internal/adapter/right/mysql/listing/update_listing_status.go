package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateListingStatus(ctx context.Context, tx *sql.Tx, listingID int64, newStatus listingmodel.ListingStatus, expectedCurrent listingmodel.ListingStatus) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listings SET status = ? WHERE id = ? AND status = ?`
	defer la.ObserveOnComplete("update", query)()

	result, err := tx.ExecContext(ctx, query, newStatus, listingID, expectedCurrent)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_status.exec_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("update listing status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_status.rows_affected_error", "error", err, "listing_id", listingID)
		return fmt.Errorf("rows affected for update listing status: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
