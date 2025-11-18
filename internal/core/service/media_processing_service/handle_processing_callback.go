package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HandleProcessingCallback wraps the payload received from Step Functions or Lambda callbacks.
func (s *mediaProcessingService) HandleProcessingCallback(ctx context.Context, input HandleProcessingCallbackInput) (HandleProcessingCallbackOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	callback := input.Callback
	if callback.JobID == 0 {
		return HandleProcessingCallbackOutput{}, derrors.Validation("jobId is required", map[string]any{"jobId": "required"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.callback.tx_start_error", "err", txErr, "job_id", callback.JobID)
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.callback.tx_rollback_error", "err", rbErr, "job_id", callback.JobID)
			}
		}
	}()

	job, err := s.repo.GetProcessingJobByID(ctx, tx, callback.JobID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Warn("service.media.callback.job_not_found", "job_id", callback.JobID)
			return HandleProcessingCallbackOutput{}, derrors.NotFound("processing job not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.get_job_error", "err", err, "job_id", callback.JobID)
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to load processing job", err)
	}

	var batchStatus mediaprocessingmodel.BatchStatus
	var jobPayload mediaprocessingmodel.MediaProcessingJobPayload
	statusMessage := "processing_completed"

	switch callback.Status {
	case mediaprocessingmodel.MediaProcessingJobStatusSucceeded:
		batchStatus = mediaprocessingmodel.BatchStatusReady
		if len(callback.Outputs) > 0 {
			jobPayload = callback.Outputs[0]
		}
	case mediaprocessingmodel.MediaProcessingJobStatusFailed:
		batchStatus = mediaprocessingmodel.BatchStatusFailed
		statusMessage = "processing_failed"
		jobPayload.ErrorMessage = callback.FailureReason
	default:
		logger.Warn("service.media.callback.unknown_status",
			"job_id", callback.JobID,
			"status", callback.Status,
		)
		return HandleProcessingCallbackOutput{}, derrors.Validation("unsupported job status", map[string]any{"status": callback.Status})
	}

	if err := s.repo.UpdateProcessingJob(ctx, tx, callback.JobID, callback.Status, jobPayload); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.update_job_error", "err", err, "job_id", callback.JobID)
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to update job", err)
	}

	details := map[string]string{
		"provider": string(callback.Provider),
		"status":   string(callback.Status),
		"jobId":    strconv.FormatUint(callback.JobID, 10),
	}
	if externalID := job.ExternalID(); externalID != "" {
		details["externalId"] = externalID
	}
	if jobPayload.RawKey != "" {
		details["rawKey"] = jobPayload.RawKey
	}
	if jobPayload.ProcessedKey != "" {
		details["processedKey"] = jobPayload.ProcessedKey
	}
	if jobPayload.ThumbnailKey != "" {
		details["thumbnailKey"] = jobPayload.ThumbnailKey
	}
	if jobPayload.ErrorMessage != "" {
		details["error"] = jobPayload.ErrorMessage
	}

	metadata := mediaprocessingmodel.BatchStatusMetadata{
		Message:   statusMessage,
		Reason:    jobPayload.ErrorMessage,
		Details:   details,
		UpdatedBy: 0,
		UpdatedAt: s.nowUTC(),
	}

	if err := s.repo.UpdateBatchStatus(ctx, tx, job.BatchID(), batchStatus, metadata); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.update_batch_status_error", "err", err, "job_id", callback.JobID, "batch_id", job.BatchID())
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to update batch status", err)
	}

	logger.Info("service.media.callback.success",
		"job_id", callback.JobID,
		"status", callback.Status,
		"batch_status", batchStatus,
		"batch_id", job.BatchID(),
		"listing_identity_id", job.ListingID(),
	)

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.tx_commit_error", "err", err, "job_id", callback.JobID)
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to commit callback", err)
	}
	committed = true

	// Acknowledge the message from the queue if receipt handle provided
	if input.ReceiptHandle != "" {
		if ackErr := s.queue.Acknowledge(ctx, input.ReceiptHandle); ackErr != nil {
			logger.Error("service.media.callback.ack_error", "err", ackErr, "job_id", callback.JobID)
			// Don't fail the entire operation if ack fails
		}
	}

	return HandleProcessingCallbackOutput{
		ListingIdentityID: listingmodel.ListingIdentityID(job.ListingID()),
		BatchID:           job.BatchID(),
		Status:            batchStatus,
	}, nil
}
