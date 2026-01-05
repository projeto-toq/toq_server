package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"time"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

const selectStuckJobsQuery = `
SELECT id, listing_identity_id, status, provider, external_id, payload, started_at, completed_at, last_error, callback_body
FROM media_processing_jobs
WHERE status = ?
  AND started_at IS NOT NULL
  AND started_at < ?
  AND completed_at IS NULL
ORDER BY started_at ASC, id ASC
LIMIT 200
`

// ListStuckJobs returns processing jobs that remained running beyond the configured timeout.
func (a *MediaProcessingAdapter) ListStuckJobs(ctx context.Context, tx *sql.Tx, status mediaprocessingmodel.MediaProcessingJobStatus, startedBefore time.Time) ([]mediaprocessingmodel.MediaProcessingJob, error) {
	observer := a.ObserveOnComplete("select", selectStuckJobsQuery)
	defer observer()

	rows, err := a.QueryContext(ctx, tx, "select", selectStuckJobsQuery, string(status), startedBefore)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]mediaprocessingmodel.MediaProcessingJob, 0)
	for rows.Next() {
		var entity mediaprocessingentities.JobEntity
		if scanErr := rows.Scan(
			&entity.ID,
			&entity.ListingIdentityID,
			&entity.Status,
			&entity.Provider,
			&entity.ExternalID,
			&entity.Payload,
			&entity.StartedAt,
			&entity.FinishedAt,
			&entity.LastError,
			&entity.CallbackBody,
		); scanErr != nil {
			return nil, scanErr
		}

		jobs = append(jobs, mediaprocessingconverters.JobEntityToDomain(entity))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}
