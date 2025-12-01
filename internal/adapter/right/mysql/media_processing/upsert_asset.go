package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const upsertAssetQuery = `
INSERT INTO media_assets (
    listing_id,
    asset_type,
    sequence,
    status,
    s3_key_raw,
    s3_key_processed,
    title,
    metadata
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    status = VALUES(status),
    s3_key_raw = VALUES(s3_key_raw),
    s3_key_processed = VALUES(s3_key_processed),
    title = VALUES(title),
    metadata = VALUES(metadata)
`

// UpsertAsset creates or updates a media asset.
func (a *MediaProcessingAdapter) UpsertAsset(ctx context.Context, tx *sql.Tx, asset mediaprocessingmodel.MediaAsset) error {
	entity := mediaprocessingconverters.AssetDomainToEntity(asset)

	_, err := a.ExecContext(ctx, tx, "upsert_asset", upsertAssetQuery,
		entity.ListingID,
		entity.AssetType,
		entity.Sequence,
		entity.Status,
		entity.S3KeyRaw,
		entity.S3KeyProcessed,
		entity.Title,
		entity.Metadata,
	)
	return err
}
