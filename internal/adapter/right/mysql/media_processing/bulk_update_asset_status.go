package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const bulkUpdateAssetStatusQuery = `
UPDATE media_assets
SET status = ?
WHERE listing_identity_id = ? AND status = ?
`

// BulkUpdateAssetStatus moves every asset for a listing from one status to another.
func (a *MediaProcessingAdapter) BulkUpdateAssetStatus(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, fromStatus, toStatus mediaprocessingmodel.MediaAssetStatus) error {
	_, err := a.ExecContext(ctx, tx, "bulk_update_asset_status", bulkUpdateAssetStatusQuery, string(toStatus), listingIdentityID, string(fromStatus))
	return err
}
