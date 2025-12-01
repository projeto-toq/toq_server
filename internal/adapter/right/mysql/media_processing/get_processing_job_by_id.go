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
SELECT id, listing_identity_id, status, provider, external_id, payload,
       started_at, completed_at
FROM media_processing_jobs
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
		&entity.ListingIdentityID,
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
