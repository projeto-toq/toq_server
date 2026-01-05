package mediaprocessingservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("service.media.callback.job_not_found", "job_id", input.JobID)
			return dto.HandleProcessingCallbackOutput{}, derrors.NotFound("processing job not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.callback.get_job_error", "err", err, "job_id", input.JobID)
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to get job", err)
	}

	if input.ListingIdentityID != 0 && job.ListingIdentityID() != input.ListingIdentityID {
		logger.Warn("service.media.callback.listing_mismatch", "job_id", input.JobID, "job_listing", job.ListingIdentityID(), "payload_listing", input.ListingIdentityID)
	}

	if input.RawPayload != "" {
		job.SetCallbackBody(input.RawPayload)
	}

	if input.ExecutionARN != "" && job.ExternalID() == "" {
		job.SetExternalID(input.ExecutionARN)
	}

	if input.StartedAt != nil {
		job.EnsureStartedAt(*input.StartedAt)
	}

	// Map status string to enum
	var jobStatus mediaprocessingmodel.MediaProcessingJobStatus
	partialSuccess := false
	switch input.Status {
	case "SUCCEEDED":
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusSucceeded
	case "PARTIAL_SUCCESS":
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusPartial
		partialSuccess = true
	case "FAILED", "PROCESSING_FAILED", "VALIDATION_FAILED", "TIMED_OUT": // Catch all failure modes
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusFailed
	default:
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusRunning
	}

	if jobStatus.IsTerminal() {
		payload := mediaprocessingmodel.MediaProcessingJobPayload{}
		if (jobStatus == mediaprocessingmodel.MediaProcessingJobStatusSucceeded || jobStatus == mediaprocessingmodel.MediaProcessingJobStatusPartial) &&
			job.Provider() == mediaprocessingmodel.MediaProcessingProviderStepFunctionsFinalization {
			job.ApplyFinalizationPayload(input.ZipBundles, input.AssetsZipped, input.ZipSizeBytes, input.UnzippedSizeBytes)
			payload = job.Payload()
		}
		job.MarkCompleted(jobStatus, payload, s.nowUTC())

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
		logger.Error("service.media.callback.update_job_error", "err", err, "job_id", input.JobID)
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
		resultStatus := strings.ToUpper(result.Status)
		resultFailed := resultStatus == "FAILED"
		if resultFailed {
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

		updatedAsset, mergeErr := s.applyProcessingResult(ctx, asset, result)
		if mergeErr != nil {
			utils.SetSpanError(ctx, mergeErr)
			logger.Warn("service.media.callback.asset_result_error", "asset_id", asset.ID(), "err", mergeErr)
			if errors.Is(mergeErr, errProcessedAssetMissingKey) {
				errorCodeHistogram["MISSING_PROCESSED_KEY"]++
			}
		}

		if updatedAsset.Status() == mediaprocessingmodel.MediaAssetStatusFailed && !resultFailed {
			failedCount++
		}

		if err := s.repo.UpsertAsset(ctx, tx, updatedAsset); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.callback.update_asset_failed", "asset_id", updatedAsset.ID(), "err", err)
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
				"partial_success", partialSuccess,
			)
		} else {
			logger.Info("service.media.callback.assets_processed",
				"job_id", input.JobID,
				"processed_assets", len(input.Results),
				"partial_success", partialSuccess,
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
