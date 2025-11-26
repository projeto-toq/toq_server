package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const listAssetsByBatchQuery = `
SELECT id, batch_id, asset_type, orientation, filename, content_type, sequence, size_bytes, checksum,
	   raw_object_key, processed_key, thumbnail_key, width, height, duration_millis, title, metadata
FROM listing_media_assets
WHERE batch_id = ?
ORDER BY sequence ASC
`

// ListAssetsByBatch retorna assets associados ao lote.
func (a *MediaProcessingAdapter) ListAssetsByBatch(ctx context.Context, tx *sql.Tx, batchID uint64) ([]mediaprocessingmodel.MediaAsset, error) {
	observer := a.ObserveOnComplete("select", listAssetsByBatchQuery)
	defer observer()

	rows, err := a.QueryContext(ctx, tx, "select", listAssetsByBatchQuery, batchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []mediaprocessingmodel.MediaAsset
	for rows.Next() {
		var entity mediaprocessingentities.AssetEntity
		if err := rows.Scan(
			&entity.ID,
			&entity.BatchID,
			&entity.AssetType,
			&entity.Orientation,
			&entity.Filename,
			&entity.ContentType,
			&entity.Sequence,
			&entity.SizeBytes,
			&entity.Checksum,
			&entity.RawObjectKey,
			&entity.ProcessedKey,
			&entity.ThumbnailKey,
			&entity.Width,
			&entity.Height,
			&entity.DurationMillis,
			&entity.Title,
			&entity.Metadata,
		); err != nil {
			return nil, err
		}
		assets = append(assets, mediaprocessingconverters.AssetEntityToDomain(entity))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}
