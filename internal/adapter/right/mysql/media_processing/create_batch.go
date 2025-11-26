package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const insertBatchQuery = `
INSERT INTO listing_media_batches (
    listing_id,
    photographer_user_id,
    status,
    upload_manifest_json,
    processing_metadata_json,
    error_code,
    error_detail,
    received_at,
    deleted_at
) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), ?)
`

// CreateBatch persiste um novo lote de mídia.
func (a *MediaProcessingAdapter) CreateBatch(ctx context.Context, tx *sql.Tx, batch mediaprocessingmodel.MediaBatch) (uint64, error) {
	entity, err := mediaprocessingconverters.BatchDomainToEntity(batch)
	if err != nil {
		return 0, err
	}

	observer := a.ObserveOnComplete("insert", insertBatchQuery)
	defer observer()

	result, err := a.ExecContext(ctx, tx, "insert", insertBatchQuery,
		entity.ListingID,
		entity.PhotographerUserID,
		entity.Status,
		entity.UploadManifestJSON,
		entity.ProcessingMetadataJSON,
		entity.ErrorCode,
		entity.ErrorDetail,
		entity.DeletedAt,
	)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(insertedID), nil
}
