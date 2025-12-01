package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const getAssetByRawKeyQuery = `
SELECT
    id, listing_identity_id, asset_type, sequence, status, s3_key_raw, s3_key_processed, title, metadata
FROM media_assets
WHERE s3_key_raw = ?
`

// GetAssetByRawKey retrieves a specific asset by its raw S3 key.
func (a *MediaProcessingAdapter) GetAssetByRawKey(ctx context.Context, tx *sql.Tx, rawKey string) (mediaprocessingmodel.MediaAsset, error) {
	var entity mediaprocessingentities.AssetEntity
	err := a.QueryRowContext(ctx, tx, "get_asset_by_raw_key", getAssetByRawKeyQuery, rawKey).Scan(
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
