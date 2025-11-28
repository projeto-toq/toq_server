package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

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

	// LOG: Raw input at service level
	logger.Info("service.media.callback.received",
		"job_id", callback.JobID,
		"status", callback.Status,
		"outputs_count", len(callback.Outputs),
	)

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

	assets, err := s.repo.ListAssetsByBatch(ctx, tx, job.BatchID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.list_assets_error", "err", err, "batch_id", job.BatchID())
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to list assets", err)
	}

	assetMap := make(map[string]*mediaprocessingmodel.MediaAsset, len(assets))
	for i := range assets {
		assetMap[assets[i].RawObjectKey()] = &assets[i]
	}

	var batchStatus mediaprocessingmodel.BatchStatus
	var failureReason string
	statusMessage := "processing_completed"
	updatedAssets := make([]mediaprocessingmodel.MediaAsset, 0, len(assets))

	switch callback.Status {
	case mediaprocessingmodel.MediaProcessingJobStatusSucceeded:
		batchStatus = mediaprocessingmodel.BatchStatusReady

		for _, output := range callback.Outputs {
			// LOG: Debugging incoming output
			logger.Debug("service.media.callback.processing_asset",
				"raw_key", output.RawKey,
				"processed_key", output.ProcessedKey,
				"thumbnail_key", output.ThumbnailKey,
			)

			if asset, exists := assetMap[output.RawKey]; exists {
				width, _ := strconv.ParseUint(output.Outputs["width"], 10, 16)
				height, _ := strconv.ParseUint(output.Outputs["height"], 10, 16)
				duration, _ := strconv.ParseUint(output.Outputs["durationMillis"], 10, 32)

				asset.SetProcessedOutputs(
					output.ProcessedKey,
					output.ThumbnailKey,
					uint16(width),
					uint16(height),
					uint32(duration),
				)

				if len(output.Outputs) > 0 {
					for k, v := range output.Outputs {
						asset.SetMetadata("variant_"+k, v)
					}
				}

				logger.Info("service.media.callback.asset_updated",
					"asset_id", asset.ID(),
					"raw_key", output.RawKey,
					"processed_key", output.ProcessedKey,
					"thumb_key", output.ThumbnailKey,
				)

				updatedAssets = append(updatedAssets, *asset)
			} else {
				if isZipArtifact(output.ProcessedKey) {
					newAsset := mediaprocessingmodel.NewMediaAsset(
						job.BatchID(),
						job.ListingID(),
						mediaprocessingmodel.MediaAssetTypeZip,
						0,
					)
					newAsset.SetProcessedOutputs(output.ProcessedKey, "", 0, 0, 0)
					newAsset.SetFilename("full_download.zip")
					newAsset.SetMetadata("title", "Pacote Completo (ZIP)")

					updatedAssets = append(updatedAssets, newAsset)
				} else {
					// LOG: Critical error diagnosis
					logger.Error("service.media.callback.asset_mismatch",
						"received_raw_key", output.RawKey,
						"batch_id", job.BatchID(),
						"available_keys_count", len(assetMap),
					)
					// Log sample keys
					keys := make([]string, 0, 3)
					for k := range assetMap {
						if len(keys) < 3 {
							keys = append(keys, k)
						}
					}
					logger.Debug("service.media.callback.available_keys_sample", "keys", keys)
				}
			}
		}

		if len(updatedAssets) > 0 {
			if err := s.repo.UpsertAssets(ctx, tx, updatedAssets); err != nil {
				utils.SetSpanError(ctx, err)
				logger.Error("service.media.callback.upsert_assets_error", "err", err, "batch_id", job.BatchID())
				return HandleProcessingCallbackOutput{}, derrors.Infra("failed to update assets", err)
			}
		}

	case mediaprocessingmodel.MediaProcessingJobStatusFailed:
		batchStatus = mediaprocessingmodel.BatchStatusFailed
		statusMessage = "processing_failed"
		failureReason = callback.FailureReason

	default:
		logger.Warn("service.media.callback.unknown_status",
			"job_id", callback.JobID,
			"status", callback.Status,
		)
		return HandleProcessingCallbackOutput{}, derrors.Validation("unsupported job status", map[string]any{"status": callback.Status})
	}

	var summaryPayload mediaprocessingmodel.MediaProcessingJobPayload
	if len(callback.Outputs) > 0 {
		summaryPayload = callback.Outputs[0]
	}
	summaryPayload.ErrorMessage = failureReason

	job.MarkCompleted(callback.Status, summaryPayload, s.nowUTC())
	if callback.ExternalID != "" {
		job.SetExternalID(callback.ExternalID)
	}

	if err := s.repo.UpdateProcessingJob(ctx, tx, job); err != nil {
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
	if summaryPayload.RawKey != "" {
		details["rawKey"] = summaryPayload.RawKey
	}
	if summaryPayload.ProcessedKey != "" {
		details["processedKey"] = summaryPayload.ProcessedKey
	}
	if summaryPayload.ThumbnailKey != "" {
		details["thumbnailKey"] = summaryPayload.ThumbnailKey
	}
	if summaryPayload.ErrorMessage != "" {
		details["error"] = summaryPayload.ErrorMessage
	}

	metadata := mediaprocessingmodel.BatchStatusMetadata{
		Message:   statusMessage,
		Reason:    failureReason,
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
		"assets_updated", len(updatedAssets),
	)

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.tx_commit_error", "err", err, "job_id", callback.JobID)
		return HandleProcessingCallbackOutput{}, derrors.Infra("failed to commit callback", err)
	}
	committed = true

	if input.ReceiptHandle != "" {
		if ackErr := s.queue.Acknowledge(ctx, input.ReceiptHandle); ackErr != nil {
			logger.Error("service.media.callback.ack_error", "err", ackErr, "job_id", callback.JobID)
		}
	}

	return HandleProcessingCallbackOutput{
		ListingIdentityID: listingmodel.ListingIdentityID(job.ListingID()),
		BatchID:           job.BatchID(),
		Status:            batchStatus,
	}, nil
}

func isZipArtifact(key string) bool {
	return strings.HasSuffix(strings.ToLower(key), ".zip")
}
