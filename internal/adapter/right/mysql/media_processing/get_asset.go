package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const getAssetQuery = `
SELECT
    id, listing_id, asset_type, sequence, status, s3_key_raw, s3_key_processed, title, metadata
FROM media_assets
WHERE listing_id = ? AND asset_type = ? AND sequence = ?
`

// GetAsset retrieves a specific asset.
func (a *MediaProcessingAdapter) GetAsset(ctx context.Context, tx *sql.Tx, listingID uint64, assetType mediaprocessingmodel.MediaAssetType, sequence uint8) (mediaprocessingmodel.MediaAsset, error) {
	var entity mediaprocessingentities.AssetEntity
	err := a.QueryRowContext(ctx, tx, "get_asset", getAssetQuery, listingID, string(assetType), sequence).Scan(
		&entity.ID,
		&entity.ListingID,
		&entity.AssetType,
		&entity.Sequence,
		&entity.Status,
		&entity.S3KeyRaw,
		&entity.S3KeyProcessed,
		&entity.Title,
		&entity.Metadata,
	)
	if err != nil {
		return mediaprocessingmodel.MediaAsset{}, err
	}

	return mediaprocessingconverters.AssetEntityToDomain(entity), nil
}
