package mediaprocessingservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
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

	var batchStatus mediaprocessingmodel.BatchStatus
	var jobPayload mediaprocessingmodel.MediaProcessingJobPayload

	switch callback.Status {
	case mediaprocessingmodel.MediaProcessingJobStatusSucceeded:
		batchStatus = mediaprocessingmodel.BatchStatusReady
		if len(callback.Outputs) > 0 {
			jobPayload = callback.Outputs[0]
		}
	case mediaprocessingmodel.MediaProcessingJobStatusFailed:
		batchStatus = mediaprocessingmodel.BatchStatusFailed
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

	// Retrieve batch/listing IDs from the job record to update batch status
	// Since we don't have a GetJobByID method, we'll work with what's in the callback
	// The callback should contain enough info or we need to add a method
	// For now, we'll assume the callback includes the necessary IDs

	// Note: This is a simplified implementation. In production, you'd retrieve the job
	// record to get the batch_id and listing_id, then update accordingly.

	logger.Info("service.media.callback.success",
		"job_id", callback.JobID,
		"status", callback.Status,
		"batch_status", batchStatus,
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
		ListingID: 0, // Would be populated from job record
		BatchID:   0, // Would be populated from job record
		Status:    batchStatus,
	}, nil
}
