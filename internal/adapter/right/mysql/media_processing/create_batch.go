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
    reference,
    status,
    status_message,
    status_reason,
    status_details,
    status_updated_by,
	status_updated_at,
	deleted_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// CreateBatch persiste um novo lote de mídia.
func (a *MediaProcessingAdapter) CreateBatch(ctx context.Context, tx *sql.Tx, batch mediaprocessingmodel.MediaBatch) (uint64, error) {
	entity := mediaprocessingconverters.BatchDomainToEntity(batch)
	observer := a.ObserveOnComplete("insert", insertBatchQuery)
	defer observer()

	result, err := a.ExecContext(ctx, tx, "insert", insertBatchQuery,
		entity.ListingID,
		entity.Reference,
		entity.Status,
		entity.StatusMessage,
		entity.StatusReason,
		entity.StatusDetails,
		entity.StatusUpdatedBy,
		entity.StatusUpdatedAt,
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
