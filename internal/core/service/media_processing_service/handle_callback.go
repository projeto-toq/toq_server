package mediaprocessingservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *mediaProcessingService) HandleProcessingCallback(ctx context.Context, input dto.HandleProcessingCallbackInput) (dto.HandleProcessingCallbackOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	logger.Info("service.media.callback.received", "job_id", input.JobID, "status", input.Status)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.callback.tx_start_error", "err", txErr, "job_id", input.JobID)
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("service.media.callback.rollback_error", "err", rbErr)
			}
		}
	}()

	// Update Job Status
	job, err := s.repo.GetProcessingJobByID(ctx, tx, input.JobID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to get job", err)
	}

	if input.ListingIdentityID != 0 && job.ListingIdentityID() != input.ListingIdentityID {
		logger.Warn("service.media.callback.listing_mismatch", "job_id", input.JobID, "job_listing", job.ListingIdentityID(), "payload_listing", input.ListingIdentityID)
	}

	if input.RawPayload != "" {
		job.SetCallbackBody(input.RawPayload)
	}

	// Map status string to enum
	var jobStatus mediaprocessingmodel.MediaProcessingJobStatus
	switch input.Status {
	case "SUCCEEDED":
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusSucceeded
	case "FAILED", "PROCESSING_FAILED", "VALIDATION_FAILED", "TIMED_OUT": // Catch all failure modes
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusFailed
	default:
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusRunning
	}

	if jobStatus.IsTerminal() {
		job.MarkCompleted(jobStatus, mediaprocessingmodel.MediaProcessingJobPayload{}, s.nowUTC())

		errorFragments := make([]string, 0, 4)
		if input.ErrorCode != "" {
			errorFragments = append(errorFragments, fmt.Sprintf("code=%s", input.ErrorCode))
		}
		if input.Error != "" {
			errorFragments = append(errorFragments, input.Error)
		}
		if input.FailureReason != "" {
			errorFragments = append(errorFragments, input.FailureReason)
		}
		if len(input.ErrorMetadata) > 0 {
			if metaBytes, marshalErr := json.Marshal(input.ErrorMetadata); marshalErr == nil {
				errorFragments = append(errorFragments, fmt.Sprintf("meta=%s", string(metaBytes)))
			}
		}
		if len(errorFragments) > 0 {
			job.AppendError(strings.Join(errorFragments, " | "))
		}
	}

	if err := s.repo.UpdateProcessingJob(ctx, tx, job); err != nil {
		utils.SetSpanError(ctx, err)
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to update job", err)
	}

	// FAIL-SAFE: If job failed globally, mark all associated assets as FAILED
	// This prevents assets from being stuck in PROCESSING forever.
	if jobStatus == mediaprocessingmodel.MediaProcessingJobStatusFailed && len(input.Results) == 0 {
		if err := s.repo.BulkUpdateAssetStatus(ctx, tx, job.ListingIdentityID(), mediaprocessingmodel.MediaAssetStatusProcessing, mediaprocessingmodel.MediaAssetStatusFailed); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.callback.bulk_fail_assets_error", "err", err, "listing_identity_id", job.ListingIdentityID())
		}
	}

	failedCount := 0
	errorCodeHistogram := make(map[string]int)

	// Update assets based on results
	for _, result := range input.Results {
		if strings.EqualFold(result.Status, "FAILED") {
			failedCount++
		}
		if result.ErrorCode != "" {
			errorCodeHistogram[result.ErrorCode]++
		}

		asset, assetErr := s.getAssetForCallback(ctx, tx, result)
		if assetErr != nil {
			logger.Error("service.media.callback.asset_lookup_error", "asset_id", result.AssetID, "raw_key", result.RawKey, "err", assetErr)
			continue
		}

		switch result.Status {
		case "PROCESSED":
			asset.SetStatus(mediaprocessingmodel.MediaAssetStatusProcessed)
			if result.ProcessedKey != "" {
				asset.SetS3KeyProcessed(result.ProcessedKey)
			}

			// Ensure metadata map exists if we have thumbnail
			if result.ThumbnailKey != "" {
				if result.Metadata == nil {
					result.Metadata = make(map[string]string)
				}
				result.Metadata["thumbnailKey"] = result.ThumbnailKey
			}

			if len(result.Metadata) > 0 {
				// Merge metadata
				currentMeta := make(map[string]string)
				if asset.Metadata() != "" {
					_ = json.Unmarshal([]byte(asset.Metadata()), &currentMeta)
				}
				for k, v := range result.Metadata {
					currentMeta[k] = v
				}
				metaBytes, _ := json.Marshal(currentMeta)
				asset.SetMetadata(string(metaBytes))
			}
		case "FAILED":
			asset.SetStatus(mediaprocessingmodel.MediaAssetStatusFailed)
			// Maybe store error in metadata?
			if result.Error != "" || result.ErrorCode != "" {
				currentMeta := make(map[string]string)
				if asset.Metadata() != "" {
					_ = json.Unmarshal([]byte(asset.Metadata()), &currentMeta)
				}
				if result.Error != "" {
					currentMeta["error"] = result.Error
				}
				if result.ErrorCode != "" {
					currentMeta["errorCode"] = result.ErrorCode
				}
				metaBytes, _ := json.Marshal(currentMeta)
				asset.SetMetadata(string(metaBytes))
			}
		}

		if err := s.repo.UpsertAsset(ctx, tx, asset); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.callback.update_asset_failed", "asset_id", asset.ID(), "err", err)
			return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to update asset", err)
		}
	}

	if len(input.Results) > 0 {
		if failedCount > 0 || input.ErrorCode != "" {
			logger.Warn("service.media.callback.assets_failed",
				"job_id", input.JobID,
				"failed_assets", failedCount,
				"error_codes", errorCodeHistogram,
				"callback_error_code", input.ErrorCode,
				"callback_error_metadata", input.ErrorMetadata,
			)
		} else {
			logger.Info("service.media.callback.assets_processed",
				"job_id", input.JobID,
				"processed_assets", len(input.Results),
			)
		}
	} else if input.ErrorCode != "" || input.Error != "" {
		logger.Warn("service.media.callback.no_results_failure",
			"job_id", input.JobID,
			"status", input.Status,
			"callback_error_code", input.ErrorCode,
			"callback_error", input.Error,
			"callback_error_metadata", input.ErrorMetadata,
		)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.tx_commit_error", "err", err, "job_id", input.JobID)
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	return dto.HandleProcessingCallbackOutput{Success: true}, nil
}

func (s *mediaProcessingService) getAssetForCallback(ctx context.Context, tx *sql.Tx, result dto.ProcessingResult) (mediaprocessingmodel.MediaAsset, error) {
	if result.AssetID != 0 {
		return s.repo.GetAssetByID(ctx, tx, result.AssetID)
	}
	if result.RawKey != "" {
		return s.repo.GetAssetByRawKey(ctx, tx, result.RawKey)
	}
	return mediaprocessingmodel.MediaAsset{}, derrors.Validation("missing asset identifier", map[string]any{"result": result})
}
