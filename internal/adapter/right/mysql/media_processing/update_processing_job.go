package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const updateProcessingJobQuery = `
UPDATE listing_media_jobs
SET status = ?, payload = ?, last_error = ?
WHERE id = ?
`

// UpdateProcessingJob atualiza status/payload do job.
func (a *MediaProcessingAdapter) UpdateProcessingJob(ctx context.Context, tx *sql.Tx, jobID uint64, status mediaprocessingmodel.MediaProcessingJobStatus, payload mediaprocessingmodel.MediaProcessingJobPayload) error {
	observer := a.ObserveOnComplete("update", updateProcessingJobQuery)
	defer observer()

	_, err := a.ExecContext(ctx, tx, "update", updateProcessingJobQuery,
		string(status),
		mediaprocessingconverters.EncodeJobPayload(payload),
		payload.ErrorMessage,
		jobID,
	)
	return err
}
