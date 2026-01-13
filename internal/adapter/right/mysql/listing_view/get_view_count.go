package mysqllistingviewadapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetCount returns the current view counter for a listing identity.
// It normalizes sql.ErrNoRows to zero to simplify service logic.
func (a *ListingViewAdapter) GetCount(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (int64, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT views FROM listing_view_metrics WHERE listing_identity_id = ?`
	row := a.QueryRowContext(ctx, tx, "select", query, listingIdentityID)

	var views int64
	if err := row.Scan(&views); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_view.get_count.scan_error", "listing_identity_id", listingIdentityID, "err", err)
		return 0, err
	}

	return views, nil
}
