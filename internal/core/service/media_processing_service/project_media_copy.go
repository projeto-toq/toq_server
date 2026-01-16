package mediaprocessingservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// copyProjectAssetToProcessed duplicates the raw object to the processed path and returns the processed key.
func (s *mediaProcessingService) copyProjectAssetToProcessed(ctx context.Context, listingID uint64, asset mediaprocessingmodel.MediaAsset) (string, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	rawKey := strings.TrimSpace(asset.S3KeyRaw())
	if rawKey == "" {
		return "", derrors.Conflict("project asset is missing raw key", derrors.WithDetails(map[string]any{"assetType": asset.AssetType(), "sequence": asset.Sequence()}))
	}

	meta, err := s.storage.ValidateObjectChecksum(ctx, rawKey, "")
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_copy.head_raw_error", "key", rawKey, "err", err, "listing_id", listingID)
		return "", derrors.Infra("failed to validate raw object", err)
	}

	content, err := s.storage.DownloadFile(ctx, rawKey)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_copy.download_raw_error", "key", rawKey, "err", err, "listing_id", listingID)
		return "", derrors.Infra("failed to download raw object", err)
	}

	processedKey := buildProjectProcessedKey(listingID, asset)

	contentType := meta.ContentType
	if strings.TrimSpace(contentType) == "" {
		if inferred := inferContentTypeFromMetadata(asset.Metadata()); inferred != "" {
			contentType = inferred
		}
	}
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	if err := s.storage.UploadFile(ctx, processedKey, content, contentType); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_copy.upload_processed_error", "key", processedKey, "err", err, "listing_id", listingID)
		return "", derrors.Infra("failed to upload processed object", err)
	}

	return processedKey, nil
}

func buildProjectProcessedKey(listingID uint64, asset mediaprocessingmodel.MediaAsset) string {
	segment := projectMediaTypeSegment(asset.AssetType())
	filename := projectFilename(asset)
	if filename == "" {
		filename = fmt.Sprintf("asset-%d-%d.bin", asset.ListingIdentityID(), asset.Sequence())
	}
	return fmt.Sprintf("%d/processed/%s/original/%s", listingID, segment, filename)
}

func projectMediaTypeSegment(assetType mediaprocessingmodel.MediaAssetType) string {
	switch assetType {
	case mediaprocessingmodel.MediaAssetTypeProjectDoc:
		return "project/doc"
	case mediaprocessingmodel.MediaAssetTypeProjectRender:
		return "project/render"
	default:
		return "misc"
	}
}

func projectFilename(asset mediaprocessingmodel.MediaAsset) string {
	rawKey := strings.TrimSpace(asset.S3KeyRaw())
	if rawKey != "" {
		_, file := path.Split(rawKey)
		if file != "" {
			return file
		}
	}

	meta := parseMetadataMap(asset.Metadata())
	filename := strings.TrimSpace(meta["filename"])
	if filename != "" {
		return filename
	}

	return ""
}

func inferContentTypeFromMetadata(rawMeta string) string {
	meta := parseMetadataMap(rawMeta)
	return strings.TrimSpace(meta["content_type"])
}

func parseMetadataMap(raw string) map[string]string {
	meta := make(map[string]string)
	if strings.TrimSpace(raw) == "" {
		return meta
	}
	_ = json.Unmarshal([]byte(raw), &meta)
	return meta
}
