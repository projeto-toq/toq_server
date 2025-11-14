package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CheckActiveListingExists verifies if a user already has an active (non-expired/non-closed) listing.
// Returns true if an active listing exists, false otherwise.
// Active listings are those with status NOT IN (StatusDraft=1, StatusClosed=13, StatusExpired=15, StatusArchived=16).
func (la *ListingAdapter) CheckActiveListingExists(ctx context.Context, tx *sql.Tx, userID int64) (bool, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT COUNT(*) 
		FROM listing_versions lv
		JOIN listing_identities li ON lv.listing_identity_id = li.id
		WHERE li.user_id = ? 
		  AND li.deleted = 0 
		  AND lv.deleted = 0 
		  AND lv.status NOT IN (1, 13, 15, 16)
	`

	var count int
	scanErr := la.QueryRowContext(ctx, tx, "select", query, userID).Scan(&count)
	if scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.listing.check_active_exists.scan_error", "error", scanErr, "user_id", userID)
		return false, fmt.Errorf("scan active listing check: %w", scanErr)
	}

	return count > 0, nil
}
