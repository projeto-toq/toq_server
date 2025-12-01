package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"errors"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const getAssetBySequenceQuery = `
SELECT
    id, listing_identity_id, asset_type, sequence, status, s3_key_raw, s3_key_processed, title, metadata
FROM media_assets
WHERE listing_identity_id = ? AND asset_type = ? AND sequence = ?
LIMIT 1
`

// GetAssetBySequence retrieves a specific asset by its business key (Listing + Type + Sequence).
func (a *MediaProcessingAdapter) GetAssetBySequence(ctx context.Context, tx *sql.Tx, listingID uint64, assetType mediaprocessingmodel.MediaAssetType, sequence uint8) (mediaprocessingmodel.MediaAsset, error) {
	var entity mediaprocessingentities.AssetEntity

	err := a.QueryRowContext(ctx, tx, "get_asset_by_sequence", getAssetBySequenceQuery, listingID, assetType, sequence).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return mediaprocessingmodel.MediaAsset{}, err // Caller handles NotFound
		}
		return mediaprocessingmodel.MediaAsset{}, err
	}

	return mediaprocessingconverters.AssetEntityToDomain(entity), nil
}
