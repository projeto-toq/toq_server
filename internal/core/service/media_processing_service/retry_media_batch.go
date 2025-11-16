package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RetryMediaBatch allows the caller to re-enqueue a finished batch.
func (s *mediaProcessingService) RetryMediaBatch(ctx context.Context, input RetryMediaBatchInput) (RetryMediaBatchOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return RetryMediaBatchOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingID == 0 {
		return RetryMediaBatchOutput{}, derrors.Validation("listingId must be greater than zero", map[string]any{"listingId": "required"})
	}
	if input.BatchID == 0 {
		return RetryMediaBatchOutput{}, derrors.Validation("batchId must be greater than zero", map[string]any{"batchId": "required"})
	}

	requestedBy, err := s.resolveRequestedBy(ctx, input.RequestedBy)
	if err != nil {
		return RetryMediaBatchOutput{}, err
	}
	input.RequestedBy = requestedBy

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.retry_batch.tx_start_error", "err", txErr, "listing_id", input.ListingID)
		return RetryMediaBatchOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.retry_batch.tx_rollback_error", "err", rbErr, "listing_id", input.ListingID)
			}
		}
	}()

	batch, err := s.repo.GetBatchByID(ctx, tx, input.BatchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RetryMediaBatchOutput{}, derrors.NotFound("batch not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.retry_batch.get_batch_error", "err", err, "batch_id", input.BatchID)
		return RetryMediaBatchOutput{}, derrors.Infra("failed to load batch", err)
	}

	if batch.ListingID() != input.ListingID {
		return RetryMediaBatchOutput{}, derrors.Conflict("batch does not belong to listing")
	}

	if !batch.Status().IsTerminal() {
		return RetryMediaBatchOutput{}, derrors.Conflict("batch is not in a terminal state")
	}

	assets, err := s.repo.ListAssetsByBatch(ctx, tx, input.BatchID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.retry_batch.list_assets_error", "err", err, "batch_id", input.BatchID)
		return RetryMediaBatchOutput{}, derrors.Infra("failed to list assets", err)
	}

	if len(assets) == 0 {
		return RetryMediaBatchOutput{}, derrors.Conflict("no assets found to retry")
	}

	objectKeys := make([]string, 0, len(assets))
	for _, asset := range assets {
		if asset.RawObjectKey() == "" {
			return RetryMediaBatchOutput{}, derrors.Conflict("raw objects missing for retry")
		}
		objectKeys = append(objectKeys, asset.RawObjectKey())
	}

	metadata := mediaprocessingmodel.BatchStatusMetadata{
		Message:   "retry_requested",
		Reason:    input.Reason,
		Details:   map[string]string{"files": fmt.Sprintf("%d", len(assets))},
		UpdatedBy: input.RequestedBy,
		UpdatedAt: s.nowUTC(),
	}

	if err := s.repo.UpdateBatchStatus(ctx, tx, input.BatchID, mediaprocessingmodel.BatchStatusProcessing, metadata); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.retry_batch.update_status_error", "err", err, "batch_id", input.BatchID)
		return RetryMediaBatchOutput{}, derrors.Infra("failed to update batch status", err)
	}

	job := mediaprocessingmodel.NewMediaProcessingJob(input.BatchID, input.ListingID, mediaprocessingmodel.MediaProcessingProviderStepFunctions)
	jobID, err := s.repo.RegisterProcessingJob(ctx, tx, job)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.retry_batch.register_job_error", "err", err, "batch_id", input.BatchID)
		return RetryMediaBatchOutput{}, derrors.Infra("failed to register processing job", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.retry_batch.tx_commit_error", "err", err, "batch_id", input.BatchID)
		return RetryMediaBatchOutput{}, derrors.Infra("failed to commit retry", err)
	}
	committed = true

	queuePayload := mediaprocessingmodel.MediaProcessingJobMessage{
		JobID:     jobID,
		BatchID:   input.BatchID,
		ListingID: input.ListingID,
		Assets:    objectKeys,
		Retry:     1,
	}

	_, err = s.queue.EnqueueRetry(ctx, queuePayload)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.retry_batch.enqueue_error",
			"err", err,
			"listing_id", input.ListingID,
			"batch_id", input.BatchID,
			"job_id", jobID,
		)
		return RetryMediaBatchOutput{}, derrors.Infra("failed to enqueue retry job", err)
	}

	logger.Info("service.media.retry_batch.success",
		"listing_id", input.ListingID,
		"batch_id", input.BatchID,
		"job_id", jobID,
		"reason", input.Reason,
	)

	return RetryMediaBatchOutput{
		ListingID: input.ListingID,
		BatchID:   input.BatchID,
		JobID:     jobID,
		Status:    mediaprocessingmodel.BatchStatusProcessing,
	}, nil
}
