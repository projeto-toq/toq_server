package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"errors"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const selectBatchByIDQuery = `
SELECT id, listing_id, photographer_user_id, status, upload_manifest_json, processing_metadata_json,
       received_at, processing_started_at, processing_finished_at, error_code, error_detail, deleted_at
FROM listing_media_batches
WHERE id = ?
`

// GetBatchByID retorna um lote específico.
func (a *MediaProcessingAdapter) GetBatchByID(ctx context.Context, tx *sql.Tx, batchID uint64) (mediaprocessingmodel.MediaBatch, error) {
	observer := a.ObserveOnComplete("select", selectBatchByIDQuery)
	defer observer()

	row := a.QueryRowContext(ctx, tx, "select", selectBatchByIDQuery, batchID)
	var entity mediaprocessingentities.BatchEntity
	if err := row.Scan(
		&entity.ID,
		&entity.ListingID,
		&entity.PhotographerUserID,
		&entity.Status,
		&entity.UploadManifestJSON,
		&entity.ProcessingMetadataJSON,
		&entity.ReceivedAt,
		&entity.ProcessingStartedAt,
		&entity.ProcessingFinishedAt,
		&entity.ErrorCode,
		&entity.ErrorDetail,
		&entity.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mediaprocessingmodel.MediaBatch{}, err
		}
		return mediaprocessingmodel.MediaBatch{}, err
	}

	return mediaprocessingconverters.BatchEntityToDomain(entity)
}
