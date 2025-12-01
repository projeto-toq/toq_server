package mediaprocessingservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GenerateDownloadURLs generates signed URLs for specific assets requested by the client.
func (s *mediaProcessingService) GenerateDownloadURLs(ctx context.Context, input dto.GenerateDownloadURLsInput) (dto.GenerateDownloadURLsOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return dto.GenerateDownloadURLsOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return dto.GenerateDownloadURLsOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	if len(input.Requests) == 0 {
		return dto.GenerateDownloadURLsOutput{}, derrors.Validation("at least one download request is required", map[string]any{"requests": "required"})
	}

	// Start transaction (read-only intent)
	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return dto.GenerateDownloadURLsOutput{}, derrors.Infra("failed to start transaction", txErr)
	}
	defer func() {
		_ = s.globalService.RollbackTransaction(ctx, tx)
	}()

	output := dto.GenerateDownloadURLsOutput{
		ListingIdentityID: input.ListingIdentityID,
		Urls:              make([]dto.DownloadURLOutput, 0, len(input.Requests)),
	}

	// For each request, verify asset exists and is processed, then generate URL
	for _, req := range input.Requests {
		// 1. Find the specific asset
		asset, err := s.repo.GetAsset(ctx, tx, uint64(input.ListingIdentityID), req.AssetType, req.Sequence)
		if err != nil {
			logger.Warn("service.media.generate_urls.asset_not_found",
				"listing_id", input.ListingIdentityID,
				"asset_type", req.AssetType,
				"sequence", req.Sequence,
				"error", err)
			continue // Skip if not found
		}

		var signedURL storageport.SignedURL

		// 2. Determine if we can serve the file
		// If resolution is "original" and status is NOT processed, we try to serve the raw file.
		if req.Resolution == "original" && asset.Status() != mediaprocessingmodel.MediaAssetStatusProcessed {
			if asset.S3KeyRaw() == "" {
				logger.Warn("service.media.generate_urls.raw_key_missing",
					"listing_id", input.ListingIdentityID,
					"asset_type", req.AssetType,
					"sequence", req.Sequence)
				continue
			}
			signedURL, err = s.storage.GenerateDownloadURL(ctx, asset.S3KeyRaw())
		} else {
			// Standard flow: must be processed
			if asset.Status() != mediaprocessingmodel.MediaAssetStatusProcessed {
				logger.Warn("service.media.generate_urls.asset_status_invalid",
					"listing_id", input.ListingIdentityID,
					"asset_type", req.AssetType,
					"sequence", req.Sequence,
					"current_status", asset.Status(),
					"expected_status", mediaprocessingmodel.MediaAssetStatusProcessed)
				continue // Skip if not processed
			}
			signedURL, err = s.storage.GenerateProcessedDownloadURL(ctx, uint64(input.ListingIdentityID), asset, req.Resolution)
		}

		if err != nil {
			logger.Error("service.media.generate_urls.generate_failed",
				"listing_id", input.ListingIdentityID,
				"asset_type", req.AssetType,
				"sequence", req.Sequence,
				"resolution", req.Resolution,
				"error", err)
			continue
		}

		output.Urls = append(output.Urls, dto.DownloadURLOutput{
			AssetType:  asset.AssetType(),
			Sequence:   asset.Sequence(),
			Resolution: req.Resolution,
			Url:        signedURL.URL,
			ExpiresIn:  int(signedURL.ExpiresIn.Seconds()),
		})
	}

	return output, nil
}
