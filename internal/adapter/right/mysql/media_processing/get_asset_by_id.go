package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const getAssetByIDQuery = `
SELECT
    id, listing_identity_id, asset_type, sequence, status, s3_key_raw, s3_key_processed, title, metadata
FROM media_assets
WHERE id = ?
`

// GetAssetByID retrieves a specific asset by its ID.
func (a *MediaProcessingAdapter) GetAssetByID(ctx context.Context, tx *sql.Tx, assetID uint64) (mediaprocessingmodel.MediaAsset, error) {
	var entity mediaprocessingentities.AssetEntity
	err := a.QueryRowContext(ctx, tx, "get_asset_by_id", getAssetByIDQuery, assetID).Scan(
		&entity.ID,
		&entity.ListingIdentityID,
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
