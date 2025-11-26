package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"encoding/json"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const updateBatchStatusQuery = `
UPDATE listing_media_batches
SET status = ?, error_code = ?, error_detail = ?, processing_metadata_json = ?
WHERE id = ? AND deleted_at IS NULL
`

// UpdateBatchStatus aplica mudanças de status e metadados.
func (a *MediaProcessingAdapter) UpdateBatchStatus(ctx context.Context, tx *sql.Tx, batchID uint64, status mediaprocessingmodel.BatchStatus, metadata mediaprocessingmodel.BatchStatusMetadata) error {
	observer := a.ObserveOnComplete("update", updateBatchStatusQuery)
	defer observer()

	var detailsJSON []byte
	if len(metadata.Details) > 0 {
		var err error
		detailsJSON, err = json.Marshal(metadata.Details)
		if err != nil {
			return err
		}
	}

	_, err := a.ExecContext(ctx, tx, "update", updateBatchStatusQuery,
		status.String(),
		metadata.Reason,  // Mapped to error_code
		metadata.Message, // Mapped to error_detail
		detailsJSON,
		batchID,
	)
	return err
}
