package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/mediaprocessingrepository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListDownloadURLs requests signed GET URLs for processed assets.
func (s *mediaProcessingService) ListDownloadURLs(ctx context.Context, input ListDownloadURLsInput) (ListDownloadURLsOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListDownloadURLsOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return ListDownloadURLsOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.list_downloads.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return ListDownloadURLsOutput{}, derrors.Infra("failed to start transaction", txErr)
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("service.media.list_downloads.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
		}
	}()

	var batch mediaprocessingmodel.MediaBatch
	if input.BatchID > 0 {
		batch, err = s.repo.GetBatchByID(ctx, tx, input.BatchID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ListDownloadURLsOutput{}, derrors.NotFound("batch not found")
			}
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.list_downloads.get_batch_error", "err", err, "batch_id", input.BatchID)
			return ListDownloadURLsOutput{}, derrors.Infra("failed to load batch", err)
		}

		if listingmodel.ListingIdentityID(batch.ListingID()) != input.ListingIdentityID {
			return ListDownloadURLsOutput{}, derrors.Conflict("batch does not belong to listing")
		}
	} else {
		filter := mediaprocessingrepository.BatchQueryFilter{
			ListingID: input.ListingIdentityID.Uint64(),
			Statuses:  []mediaprocessingmodel.BatchStatus{mediaprocessingmodel.BatchStatusReady},
			Limit:     1,
		}
		batches, err := s.repo.ListBatchesByListing(ctx, tx, filter)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.list_downloads.list_batches_error", "err", err, "listing_identity_id", input.ListingIdentityID)
			return ListDownloadURLsOutput{}, derrors.Infra("failed to list batches", err)
		}
		if len(batches) == 0 {
			return ListDownloadURLsOutput{}, derrors.NotFound("no ready batch found for listing")
		}
		batch = batches[0]
	}

	if batch.Status() != mediaprocessingmodel.BatchStatusReady {
		return ListDownloadURLsOutput{}, derrors.Conflict("batch is not ready for downloads")
	}

	assets, err := s.repo.ListAssetsByBatch(ctx, tx, batch.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.list_downloads.list_assets_error", "err", err, "batch_id", batch.ID())
		return ListDownloadURLsOutput{}, derrors.Infra("failed to list assets", err)
	}

	if len(assets) == 0 {
		return ListDownloadURLsOutput{}, derrors.NotFound("no assets found in batch")
	}

	downloads := make([]DownloadEntry, 0, len(assets))
	generatedAt := s.nowUTC()
	var maxTTL time.Duration

	for _, asset := range assets {
		if asset.ProcessedKey() == "" {
			logger.Warn("service.media.list_downloads.missing_processed_key",
				"listing_identity_id", input.ListingIdentityID,
				"batch_id", batch.ID(),
				"asset_id", asset.ID(),
			)
			continue
		}

		signedURL, err := s.storage.GenerateProcessedDownloadURL(ctx, input.ListingIdentityID.Uint64(), asset)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.list_downloads.generate_url_error",
				"err", err,
				"listing_identity_id", input.ListingIdentityID,
				"asset_id", asset.ID(),
			)
			return ListDownloadURLsOutput{}, derrors.Infra("failed to generate download URL", err)
		}

		if signedURL.ExpiresIn > maxTTL {
			maxTTL = signedURL.ExpiresIn
		}

		previewURL := ""
		if asset.ThumbnailKey() != "" {
			previewAsset := asset
			previewSignedURL, previewErr := s.storage.GenerateProcessedDownloadURL(ctx, input.ListingIdentityID.Uint64(), previewAsset)
			if previewErr == nil {
				previewURL = previewSignedURL.URL
			}
		}

		clientID := asset.Metadata()["client_id"]
		title := asset.Metadata()["title"]

		downloads = append(downloads, DownloadEntry{
			ClientID:   clientID,
			AssetType:  asset.AssetType(),
			Title:      title,
			Sequence:   asset.Sequence(),
			URL:        signedURL.URL,
			ExpiresAt:  generatedAt.Add(signedURL.ExpiresIn),
			PreviewURL: previewURL,
			Metadata:   cloneStringMap(asset.Metadata()),
		})
	}

	if len(downloads) == 0 {
		return ListDownloadURLsOutput{}, derrors.NotFound("no processed assets available")
	}

	logger.Info("service.media.list_downloads.success",
		"listing_identity_id", input.ListingIdentityID,
		"batch_id", batch.ID(),
		"downloads", len(downloads),
	)

	return ListDownloadURLsOutput{
		ListingIdentityID: input.ListingIdentityID,
		BatchID:           batch.ID(),
		GeneratedAt:       generatedAt,
		TTLSeconds:        int(maxTTL.Seconds()),
		Downloads:         downloads,
	}, nil
}
