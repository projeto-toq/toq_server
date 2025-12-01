package mediaprocessingservice

import (
	"context"
	"encoding/json"

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

	logger := utils.LoggerFromContext(ctx)
	logger.Info("service.media.callback.received", "job_id", input.JobID, "status", input.Status)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
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
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to get job", err)
	}

	// Map status string to enum
	var jobStatus mediaprocessingmodel.MediaProcessingJobStatus
	switch input.Status {
	case "SUCCEEDED":
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusSucceeded
	case "FAILED":
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusFailed
	default:
		jobStatus = mediaprocessingmodel.MediaProcessingJobStatusRunning
	}

	if jobStatus.IsTerminal() {
		job.MarkCompleted(jobStatus, mediaprocessingmodel.MediaProcessingJobPayload{}, s.nowUTC())
	}

	if err := s.repo.UpdateProcessingJob(ctx, tx, job); err != nil {
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to update job", err)
	}

	// Update assets based on results
	for _, result := range input.Results {
		var asset mediaprocessingmodel.MediaAsset
		var err error
		if result.AssetID != 0 {
			asset, err = s.repo.GetAssetByID(ctx, tx, result.AssetID)
		} else if result.RawKey != "" {
			asset, err = s.repo.GetAssetByRawKey(ctx, tx, result.RawKey)
		} else {
			logger.Error("service.media.callback.missing_identifier", "result", result)
			continue
		}

		if err != nil {
			logger.Error("service.media.callback.asset_not_found", "asset_id", result.AssetID, "raw_key", result.RawKey, "err", err)
			continue // Skip this asset but try others
		}

		if result.Status == "PROCESSED" {
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
		} else if result.Status == "FAILED" {
			asset.SetStatus(mediaprocessingmodel.MediaAssetStatusFailed)
			// Maybe store error in metadata?
			if result.Error != "" {
				currentMeta := make(map[string]string)
				if asset.Metadata() != "" {
					_ = json.Unmarshal([]byte(asset.Metadata()), &currentMeta)
				}
				currentMeta["error"] = result.Error
				metaBytes, _ := json.Marshal(currentMeta)
				asset.SetMetadata(string(metaBytes))
			}
		}

		if err := s.repo.UpsertAsset(ctx, tx, asset); err != nil {
			logger.Error("service.media.callback.update_asset_failed", "asset_id", asset.ID(), "err", err)
			return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to update asset", err)
		}
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		return dto.HandleProcessingCallbackOutput{}, derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	return dto.HandleProcessingCallbackOutput{Success: true}, nil
}
