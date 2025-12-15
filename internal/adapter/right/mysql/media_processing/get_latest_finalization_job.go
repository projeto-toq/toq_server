package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"errors"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const selectLatestFinalizationJobQuery = `
SELECT id, listing_identity_id, status, provider, external_id, payload,
       started_at, completed_at
FROM media_processing_jobs
WHERE listing_identity_id = ?
  AND provider = ?
  AND status = ?
  AND completed_at IS NOT NULL
ORDER BY completed_at DESC, id DESC
LIMIT 1
`

// GetLatestFinalizationJob retrieves the most recent successful finalization job for a listing.
func (a *MediaProcessingAdapter) GetLatestFinalizationJob(ctx context.Context, tx *sql.Tx, listingIdentityID uint64) (mediaprocessingmodel.MediaProcessingJob, error) {
	observer := a.ObserveOnComplete("select", selectLatestFinalizationJobQuery)
	defer observer()

	row := a.QueryRowContext(ctx, tx, "select", selectLatestFinalizationJobQuery,
		listingIdentityID,
		string(mediaprocessingmodel.MediaProcessingProviderStepFunctionsFinalization),
		string(mediaprocessingmodel.MediaProcessingJobStatusSucceeded),
	)

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
