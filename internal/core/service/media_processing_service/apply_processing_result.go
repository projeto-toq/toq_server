package mediaprocessingservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

var errProcessedAssetMissingKey = errors.New("processed asset missing finalized key")

func (s *mediaProcessingService) applyProcessingResult(ctx context.Context, asset mediaprocessingmodel.MediaAsset, result dto.ProcessingResult) (mediaprocessingmodel.MediaAsset, error) {
	logger := utils.LoggerFromContext(ctx)
	status := strings.ToUpper(strings.TrimSpace(result.Status))

	switch status {
	case "PROCESSED":
		if result.ProcessedKey == "" {
			err := fmt.Errorf("%w: asset_id=%d raw_key=%s", errProcessedAssetMissingKey, asset.ID(), result.RawKey)
			logger.Warn("service.media.callback.processed_missing_key", "asset_id", asset.ID(), "raw_key", result.RawKey)
			asset.SetStatus(mediaprocessingmodel.MediaAssetStatusFailed)

			failureMeta := cloneMetadata(result.Metadata)
			failureMeta["errorCode"] = "MISSING_PROCESSED_KEY"
			failureMeta["error"] = "processed asset returned without processedKey"
			asset = mergeAssetMetadata(asset, failureMeta)
			return asset, err
		}

		asset.SetStatus(mediaprocessingmodel.MediaAssetStatusProcessed)
		asset.SetS3KeyProcessed(result.ProcessedKey)

		successMeta := cloneMetadata(result.Metadata)
		if result.ThumbnailKey != "" {
			successMeta["thumbnailKey"] = result.ThumbnailKey
		}
		asset = mergeAssetMetadata(asset, successMeta)
		return asset, nil
	default:
		asset.SetStatus(mediaprocessingmodel.MediaAssetStatusFailed)
		failureMeta := cloneMetadata(result.Metadata)
		if result.Error != "" {
			failureMeta["error"] = result.Error
		}
		if result.ErrorCode != "" {
			failureMeta["errorCode"] = result.ErrorCode
		}
		if result.ThumbnailKey != "" {
			failureMeta["thumbnailKey"] = result.ThumbnailKey
		}
		asset = mergeAssetMetadata(asset, failureMeta)
		return asset, nil
	}
}

func mergeAssetMetadata(asset mediaprocessingmodel.MediaAsset, updates map[string]string) mediaprocessingmodel.MediaAsset {
	if len(updates) == 0 {
		return asset
	}

	merged := make(map[string]string)
	if asset.Metadata() != "" {
		_ = json.Unmarshal([]byte(asset.Metadata()), &merged)
	}

	for k, v := range updates {
		if k == "" || v == "" {
			continue
		}
		merged[k] = v
	}

	if len(merged) == 0 {
		return asset
	}

	if payload, err := json.Marshal(merged); err == nil {
		asset.SetMetadata(string(payload))
	}
	return asset
}

func cloneMetadata(source map[string]string) map[string]string {
	if len(source) == 0 {
		return make(map[string]string)
	}

	clone := make(map[string]string, len(source))
	for k, v := range source {
		if k == "" || v == "" {
			continue
		}
		clone[k] = v
	}
	return clone
}
