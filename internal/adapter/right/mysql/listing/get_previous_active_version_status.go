package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetPreviousActiveVersionStatus retrieves the status of the currently active version
// for a listing identity. Used when promoting a new version to inherit the previous status.
// Returns sql.ErrNoRows if identity has no active version.
func (la *ListingAdapter) GetPreviousActiveVersionStatus(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (listingmodel.ListingStatus, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT lv.status 
		FROM listing_versions lv
		JOIN listing_identities li ON lv.id = li.active_version_id
		WHERE li.id = ? AND li.deleted = 0 AND lv.deleted = 0
		LIMIT 1
	`

	var status uint8
	scanErr := la.QueryRowContext(ctx, tx, "select", query, listingIdentityID).Scan(&status)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return 0, scanErr
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.listing.get_previous_status.scan_error", "error", scanErr, "listing_identity_id", listingIdentityID)
		return 0, fmt.Errorf("scan previous active version status: %w", scanErr)
	}

	return listingmodel.ListingStatus(status), nil
}
