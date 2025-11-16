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
SELECT id, listing_id, reference, status, status_message, status_reason, status_details, status_updated_by, status_updated_at,
	   deleted_at
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
		&entity.Reference,
		&entity.Status,
		&entity.StatusMessage,
		&entity.StatusReason,
		&entity.StatusDetails,
		&entity.StatusUpdatedBy,
		&entity.StatusUpdatedAt,
		&entity.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mediaprocessingmodel.MediaBatch{}, err
		}
		return mediaprocessingmodel.MediaBatch{}, err
	}

	return mediaprocessingconverters.BatchEntityToDomain(entity), nil
}
