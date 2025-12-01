package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const deleteAssetQuery = `
DELETE FROM media_assets
WHERE listing_id = ? AND asset_type = ? AND sequence = ?
`

// DeleteAsset removes an asset from the database.
func (a *MediaProcessingAdapter) DeleteAsset(ctx context.Context, tx *sql.Tx, listingID uint64, assetType mediaprocessingmodel.MediaAssetType, sequence uint8) error {
	_, err := a.ExecContext(ctx, tx, "delete_asset", deleteAssetQuery, listingID, string(assetType), sequence)
	return err
}
