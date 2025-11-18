package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"errors"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const selectProcessingJobByIDQuery = `
SELECT id, batch_id, listing_id, status, provider, external_id, payload, retry_count,
       started_at, completed_at, last_error, callback_raw
FROM listing_media_jobs
WHERE id = ?
`

// GetProcessingJobByID retorna os metadados de um job específico.
func (a *MediaProcessingAdapter) GetProcessingJobByID(ctx context.Context, tx *sql.Tx, jobID uint64) (mediaprocessingmodel.MediaProcessingJob, error) {
	observer := a.ObserveOnComplete("select", selectProcessingJobByIDQuery)
	defer observer()

	row := a.QueryRowContext(ctx, tx, "select", selectProcessingJobByIDQuery, jobID)
	var entity mediaprocessingentities.JobEntity
	if err := row.Scan(
		&entity.ID,
		&entity.BatchID,
		&entity.ListingID,
		&entity.Status,
		&entity.Provider,
		&entity.ExternalID,
		&entity.Payload,
		&entity.RetryCount,
		&entity.StartedAt,
		&entity.CompletedAt,
		&entity.LastError,
		&entity.CallbackRaw,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mediaprocessingmodel.MediaProcessingJob{}, err
		}
		return mediaprocessingmodel.MediaProcessingJob{}, err
	}

	return mediaprocessingconverters.JobEntityToDomain(entity), nil
}
