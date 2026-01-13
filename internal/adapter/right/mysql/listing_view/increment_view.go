package mysqllistingviewadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// IncrementAndGet atomically increments the view counter and returns the updated total.
// It leverages INSERT ... ON DUPLICATE KEY UPDATE with LAST_INSERT_ID to fetch the new value
// in a single round trip, ensuring correctness under concurrent increments.
func (a *ListingViewAdapter) IncrementAndGet(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (int64, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO listing_view_metrics (listing_identity_id, views, last_view_at)
              VALUES (?, 1, NOW())
              ON DUPLICATE KEY UPDATE views = LAST_INSERT_ID(views + 1), last_view_at = NOW()`

	if _, err := a.ExecContext(ctx, tx, "upsert", query, listingIdentityID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_view.increment.exec_error", "listing_identity_id", listingIdentityID, "err", err)
		return 0, err
	}

	row := a.QueryRowContext(ctx, tx, "select", `SELECT LAST_INSERT_ID()`)
	var total int64
	if err := row.Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_view.increment.scan_error", "listing_identity_id", listingIdentityID, "err", err)
		return 0, err
	}

	slog.Debug("mysql.listing_view.increment.ok", "listing_identity_id", listingIdentityID, "views", total)
	return total, nil
}
