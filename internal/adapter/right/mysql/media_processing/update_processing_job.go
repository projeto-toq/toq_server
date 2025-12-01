package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const updateProcessingJobQuery = `
UPDATE media_jobs
SET 
    status = ?, 
    external_job_id = ?,
    output_payload_json = ?,
    started_at = ?,
    finished_at = ?
WHERE id = ?
`

// UpdateProcessingJob updates the full state of a processing job.
//
// This method persists critical lifecycle data (external ID, start/finish times)
// that receives from async callbacks.
func (a *MediaProcessingAdapter) UpdateProcessingJob(ctx context.Context, tx *sql.Tx, job mediaprocessingmodel.MediaProcessingJob) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := mediaprocessingconverters.JobDomainToEntity(job)

	result, err := a.ExecContext(ctx, tx, "update", updateProcessingJobQuery,
		entity.Status,
		entity.ExternalID,
		entity.Payload,
		entity.StartedAt,
		entity.FinishedAt,
		entity.ID,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.media_job.update.exec_error", "err", err, "job_id", job.ID())
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
