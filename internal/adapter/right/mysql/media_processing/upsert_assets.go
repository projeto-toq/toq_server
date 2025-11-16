package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"errors"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const deleteAssetsByBatchQuery = `DELETE FROM listing_media_assets WHERE batch_id = ?`

const insertAssetQuery = `
INSERT INTO listing_media_assets (
    batch_id,
    listing_id,
    asset_type,
    orientation,
    filename,
    content_type,
    sequence,
    size_bytes,
    checksum,
    raw_object_key,
    processed_key,
    thumbnail_key,
    width,
    height,
    duration_millis,
	metadata
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// UpsertAssets substitui os assets de um lote.
func (a *MediaProcessingAdapter) UpsertAssets(ctx context.Context, tx *sql.Tx, assets []mediaprocessingmodel.MediaAsset) error {
	if len(assets) == 0 {
		return nil
	}

	batchID := assets[0].BatchID()
	for _, asset := range assets {
		if asset.BatchID() != batchID {
			return errors.New("mixed batch ids provided")
		}
	}

	if _, err := a.ExecContext(ctx, tx, "delete", deleteAssetsByBatchQuery, batchID); err != nil {
		return err
	}

	for _, asset := range assets {
		entity := mediaprocessingconverters.AssetDomainToEntity(asset)
		observer := a.ObserveOnComplete("insert", insertAssetQuery)
		if _, err := a.ExecContext(ctx, tx, "insert", insertAssetQuery,
			entity.BatchID,
			entity.ListingID,
			entity.AssetType,
			entity.Orientation,
			entity.Filename,
			entity.ContentType,
			entity.Sequence,
			entity.SizeInBytes,
			entity.Checksum,
			entity.RawObjectKey,
			entity.ProcessedKey,
			entity.ThumbnailKey,
			entity.Width,
			entity.Height,
			entity.DurationMillis,
			entity.Metadata,
		); err != nil {
			observer()
			return err
		}
		observer()
	}

	return nil
}
