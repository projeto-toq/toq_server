package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const insertProcessingJobQuery = `
INSERT INTO listing_media_jobs (
    batch_id,
    listing_id,
    status,
    provider,
    external_id,
    payload,
    retry_count,
    started_at,
    completed_at,
    last_error,
    callback_raw
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// RegisterProcessingJob cria um novo job associado ao lote.
func (a *MediaProcessingAdapter) RegisterProcessingJob(ctx context.Context, tx *sql.Tx, job mediaprocessingmodel.MediaProcessingJob) (uint64, error) {
	entity := mediaprocessingconverters.JobDomainToEntity(job)
	observer := a.ObserveOnComplete("insert", insertProcessingJobQuery)
	defer observer()

	result, err := a.ExecContext(ctx, tx, "insert", insertProcessingJobQuery,
		entity.BatchID,
		entity.ListingID,
		entity.Status,
		entity.Provider,
		entity.ExternalID,
		entity.Payload,
		entity.RetryCount,
		entity.StartedAt,
		entity.CompletedAt,
		entity.LastError,
		entity.CallbackRaw,
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
