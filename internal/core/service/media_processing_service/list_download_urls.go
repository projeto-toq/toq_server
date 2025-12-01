package mediaprocessingservice

import (
	"context"
	"encoding/json"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/mediaprocessingrepository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListDownloadURLs retrieves the list of media assets for a listing, including signed URLs for processed items.
func (s *mediaProcessingService) ListDownloadURLs(ctx context.Context, input dto.ListDownloadURLsInput) (dto.ListDownloadURLsOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return dto.ListDownloadURLsOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return dto.ListDownloadURLsOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	// Start transaction (read-only intent)
	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return dto.ListDownloadURLsOutput{}, derrors.Infra("failed to start transaction", txErr)
	}
	defer func() {
		_ = s.globalService.RollbackTransaction(ctx, tx)
	}()

	// Filter setup
	filter := mediaprocessingrepository.AssetFilter{}
	if len(input.AssetTypes) > 0 {
		filter.AssetTypes = input.AssetTypes
	}

	// Fetch assets
	assets, err := s.repo.ListAssets(ctx, tx, uint64(input.ListingIdentityID), filter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return dto.ListDownloadURLsOutput{}, derrors.Infra("failed to list assets", err)
	}

	output := dto.ListDownloadURLsOutput{
		ListingIdentityID: input.ListingIdentityID,
		Assets:            make([]dto.DownloadAsset, 0, len(assets)),
	}

	for _, asset := range assets {
		downloadAsset := dto.DownloadAsset{
			AssetType: asset.AssetType(),
			Sequence:  asset.Sequence(),
			Status:    asset.Status(),
			Title:     asset.Title(),
		}

		// Parse metadata
		if metaStr := asset.Metadata(); metaStr != "" {
			var metaMap map[string]string
			if err := json.Unmarshal([]byte(metaStr), &metaMap); err == nil {
				downloadAsset.Metadata = metaMap
			}
		}

		// Generate URL if processed
		if asset.Status() == mediaprocessingmodel.MediaAssetStatusProcessed {
			signedURL, err := s.storage.GenerateProcessedDownloadURL(ctx, uint64(input.ListingIdentityID), asset)
			if err != nil {
				// Log error but continue, partial success is better than total failure for lists
				logger.Warn("service.media.list_urls.generate_url_failed", "asset_id", asset.ID(), "err", err)
			} else {
				downloadAsset.DownloadURL = signedURL.URL
			}
		}

		output.Assets = append(output.Assets, downloadAsset)
	}

	return output, nil
}
