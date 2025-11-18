package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CompleteUploadBatch confirms that all uploads finished and enqueues the batch for async processing.
func (s *mediaProcessingService) CompleteUploadBatch(ctx context.Context, input CompleteUploadBatchInput) (CompleteUploadBatchOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return CompleteUploadBatchOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}
	if input.BatchID == 0 {
		return CompleteUploadBatchOutput{}, derrors.Validation("batchId must be greater than zero", map[string]any{"batchId": "required"})
	}

	requestedBy, err := s.resolveRequestedBy(ctx, input.RequestedBy)
	if err != nil {
		return CompleteUploadBatchOutput{}, err
	}
	input.RequestedBy = requestedBy

	if len(input.Files) == 0 {
		return CompleteUploadBatchOutput{}, derrors.Validation("files are required", map[string]any{"files": "min=1"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.complete_batch.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.complete_batch.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	batch, err := s.repo.GetBatchByID(ctx, tx, input.BatchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CompleteUploadBatchOutput{}, derrors.NotFound("batch not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.get_batch_error", "err", err, "batch_id", input.BatchID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to load batch", err)
	}

	if listingmodel.ListingIdentityID(batch.ListingID()) != input.ListingIdentityID {
		return CompleteUploadBatchOutput{}, derrors.Conflict("batch does not belong to listing")
	}

	if batch.Status() != mediaprocessingmodel.BatchStatusPendingUpload {
		return CompleteUploadBatchOutput{}, derrors.Conflict("batch is not pending upload")
	}

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID.Int64())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CompleteUploadBatchOutput{}, derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Status() != listingmodel.StatusPendingPhotoProcessing {
		return CompleteUploadBatchOutput{}, derrors.Conflict("listing is not awaiting media uploads")
	}

	existingAssets, err := s.repo.ListAssetsByBatch(ctx, tx, input.BatchID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.list_assets_error", "err", err, "batch_id", input.BatchID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to list assets", err)
	}

	if len(existingAssets) != len(input.Files) {
		return CompleteUploadBatchOutput{}, derrors.Validation("file count mismatch", map[string]any{
			"expected": len(existingAssets),
			"provided": len(input.Files),
		})
	}

	assetMap := make(map[string]mediaprocessingmodel.MediaAsset, len(existingAssets))
	for _, asset := range existingAssets {
		clientID := asset.Metadata()["client_id"]
		if clientID != "" {
			assetMap[clientID] = asset
		}
	}

	updatedAssets := make([]mediaprocessingmodel.MediaAsset, 0, len(input.Files))
	objectKeys := make([]string, 0, len(input.Files))

	for idx, file := range input.Files {
		clientID := strings.TrimSpace(file.ClientID)
		if clientID == "" {
			return CompleteUploadBatchOutput{}, derrors.Validation("clientId is required", map[string]any{
				fmt.Sprintf("files[%d].clientId", idx): "required",
			})
		}

		asset, exists := assetMap[clientID]
		if !exists {
			return CompleteUploadBatchOutput{}, derrors.Validation("unknown clientId", map[string]any{"clientId": clientID})
		}

		objectKey := strings.TrimSpace(file.ObjectKey)
		if objectKey == "" {
			return CompleteUploadBatchOutput{}, derrors.Validation("objectKey is required", map[string]any{"clientId": clientID})
		}

		if asset.RawObjectKey() != objectKey {
			return CompleteUploadBatchOutput{}, derrors.Validation("objectKey mismatch", map[string]any{
				"clientId": clientID,
				"expected": asset.RawObjectKey(),
				"provided": objectKey,
			})
		}

		metadata, err := s.storage.ValidateObjectChecksum(ctx, objectKey, asset.Checksum())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.complete_batch.validate_checksum_error",
				"err", err,
				"client_id", clientID,
				"object_key", objectKey,
			)
			return CompleteUploadBatchOutput{}, derrors.Validation("checksum validation failed", map[string]any{
				"clientId": clientID,
				"error":    err.Error(),
			})
		}

		asset.UpdateRawObject(objectKey, asset.Checksum(), metadata.SizeInBytes)
		if file.ETag != "" {
			asset.SetMetadata("etag", file.ETag)
		}

		updatedAssets = append(updatedAssets, asset)
		objectKeys = append(objectKeys, objectKey)
	}

	if err := s.repo.UpsertAssets(ctx, tx, updatedAssets); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.upsert_assets_error", "err", err, "batch_id", input.BatchID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to update assets", err)
	}

	metadata := mediaprocessingmodel.BatchStatusMetadata{
		Message:   "uploads_confirmed",
		Details:   map[string]string{"files": fmt.Sprintf("%d", len(updatedAssets))},
		UpdatedBy: input.RequestedBy,
		UpdatedAt: s.nowUTC(),
	}

	if err := s.repo.UpdateBatchStatus(ctx, tx, input.BatchID, mediaprocessingmodel.BatchStatusReceived, metadata); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.update_status_error", "err", err, "batch_id", input.BatchID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to update batch status", err)
	}

	job := mediaprocessingmodel.NewMediaProcessingJob(input.BatchID, input.ListingIdentityID.Uint64(), mediaprocessingmodel.MediaProcessingProviderStepFunctions)
	jobID, err := s.repo.RegisterProcessingJob(ctx, tx, job)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.register_job_error", "err", err, "batch_id", input.BatchID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to register processing job", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.tx_commit_error", "err", err, "batch_id", input.BatchID)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to commit batch completion", err)
	}
	committed = true

	queuePayload := mediaprocessingmodel.MediaProcessingJobMessage{
		JobID:     jobID,
		BatchID:   input.BatchID,
		ListingID: input.ListingIdentityID.Uint64(),
		Assets:    objectKeys,
		Retry:     0,
	}

	_, err = s.queue.EnqueueJob(ctx, queuePayload)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete_batch.enqueue_error",
			"err", err,
			"listing_identity_id", input.ListingIdentityID,
			"batch_id", input.BatchID,
			"job_id", jobID,
		)
		return CompleteUploadBatchOutput{}, derrors.Infra("failed to enqueue processing job", err)
	}

	logger.Info("service.media.complete_batch.success",
		"listing_identity_id", input.ListingIdentityID,
		"batch_id", input.BatchID,
		"job_id", jobID,
		"files", len(updatedAssets),
	)

	return CompleteUploadBatchOutput{
		ListingIdentityID: input.ListingIdentityID,
		BatchID:           input.BatchID,
		JobID:             jobID,
		Status:            mediaprocessingmodel.BatchStatusReceived,
		EstimatedDuration: 5 * time.Minute,
	}, nil
}
