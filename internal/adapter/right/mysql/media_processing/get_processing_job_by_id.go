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
SELECT j.id, j.batch_id, b.listing_id, j.status, j.provider, j.external_job_id, j.output_payload_json,
       j.started_at, j.finished_at
FROM listing_media_jobs j
JOIN listing_media_batches b ON j.batch_id = b.id
WHERE j.id = ?
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
		&entity.StartedAt,
		&entity.FinishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mediaprocessingmodel.MediaProcessingJob{}, err
		}
		return mediaprocessingmodel.MediaProcessingJob{}, err
	}

	return mediaprocessingconverters.JobEntityToDomain(entity), nil
}
