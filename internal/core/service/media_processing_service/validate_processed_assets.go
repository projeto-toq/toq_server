package mediaprocessingservice

import (
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// ensureAssetsReadyForFinalization validates that a listing has processed assets ready for ZIP generation.
func ensureAssetsReadyForFinalization(assets []mediaprocessingmodel.MediaAsset) ([]mediaprocessingmodel.MediaAsset, error) {
	if len(assets) == 0 {
		return nil, derrors.Conflict("no assets found for this listing")
	}

	processed := make([]mediaprocessingmodel.MediaAsset, 0, len(assets))
	for _, asset := range assets {
		switch asset.Status() {
		case mediaprocessingmodel.MediaAssetStatusProcessing:
			return nil, derrors.Conflict("assets are still processing, please wait")
		case mediaprocessingmodel.MediaAssetStatusPendingUpload:
			return nil, derrors.Conflict("assets are pending upload, please call /process endpoint first")
		case mediaprocessingmodel.MediaAssetStatusFailed:
			return nil, derrors.Conflict("some assets failed processing, please remove or retry them")
		case mediaprocessingmodel.MediaAssetStatusProcessed:
			if asset.S3KeyRaw() == "" {
				field := fmt.Sprintf("assetId:%d", asset.ID())
				return nil, derrors.Conflict(
					"processed assets are missing original keys",
					derrors.WithDetails(map[string]any{"asset": field}),
				)
			}
			if asset.S3KeyProcessed() == "" {
				field := fmt.Sprintf("assetId:%d", asset.ID())
				return nil, derrors.Conflict(
					"processed assets are missing finalized keys",
					derrors.WithDetails(map[string]any{"asset": field}),
				)
			}
			processed = append(processed, asset)
		}
	}

	if len(processed) == 0 {
		return nil, derrors.Conflict("no processed assets found to finalize")
	}

	return processed, nil
}
