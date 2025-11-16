package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const updateBatchStatusQuery = `
UPDATE listing_media_batches
SET status = ?, status_message = ?, status_reason = ?, status_details = ?, status_updated_by = ?, status_updated_at = ?
WHERE id = ? AND deleted_at IS NULL
`

// UpdateBatchStatus aplica mudanças de status e metadados.
func (a *MediaProcessingAdapter) UpdateBatchStatus(ctx context.Context, tx *sql.Tx, batchID uint64, status mediaprocessingmodel.BatchStatus, metadata mediaprocessingmodel.BatchStatusMetadata) error {
	observer := a.ObserveOnComplete("update", updateBatchStatusQuery)
	defer observer()

	_, err := a.ExecContext(ctx, tx, "update", updateBatchStatusQuery,
		status.String(),
		metadata.Message,
		metadata.Reason,
		mediaprocessingconverters.EncodeStatusDetails(metadata.Details),
		metadata.UpdatedBy,
		metadata.UpdatedAt,
		batchID,
	)
	return err
}
