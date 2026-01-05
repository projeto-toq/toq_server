package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/media_processing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ProcessMedia triggers the async processing pipeline for pending assets.
func (s *mediaProcessingService) ProcessMedia(ctx context.Context, input dto.ProcessMediaInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.process.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.process.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to load listing", err)
	}

	if listing.Status() != listingmodel.StatusPendingPhotoProcessing && listing.Status() != listingmodel.StatusRejectedByOwner {
		return derrors.Conflict("listing is not awaiting media processing")
	}

	// Find assets that need processing (new uploads or failed attempts)
	filter := mediaprocessingrepository.AssetFilter{
		Status: []mediaprocessingmodel.MediaAssetStatus{
			mediaprocessingmodel.MediaAssetStatusPendingUpload,
			mediaprocessingmodel.MediaAssetStatusFailed,
		},
	}
	assets, err := s.repo.ListAssets(ctx, tx, uint64(input.ListingIdentityID), filter, nil)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to list pending assets", err)
	}

	if len(assets) == 0 {
		return derrors.Validation("no assets available for processing", map[string]any{"listingIdentityId": input.ListingIdentityID})
	}

	if err := s.ensureRawObjectsExist(ctx, assets); err != nil {
		return err
	}

	// Register Job first to get ID
	job := mediaprocessingmodel.NewMediaProcessingJob(uint64(input.ListingIdentityID), mediaprocessingmodel.MediaProcessingProviderStepFunctions)
	jobID, err := s.repo.RegisterProcessingJob(ctx, tx, job)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to register job", err)
	}

	// Prepare payload for Step Function
	jobMsg := mediaprocessingmodel.MediaProcessingJobMessage{
		JobID:             jobID,
		ListingIdentityID: uint64(input.ListingIdentityID),
		Assets:            make([]mediaprocessingmodel.JobAsset, 0, len(assets)),
	}

	for _, asset := range assets {
		if asset.S3KeyRaw() == "" {
			logger.Warn("service.media.process.asset_missing_raw_key", "asset_id", asset.ID(), "listing_identity_id", input.ListingIdentityID)
			continue
		}

		jobMsg.Assets = append(jobMsg.Assets, mediaprocessingmodel.JobAsset{
			Key:  asset.S3KeyRaw(),
			Type: string(asset.AssetType()),
		})

		// Update status to PROCESSING to prevent double submission
		asset.SetStatus(mediaprocessingmodel.MediaAssetStatusProcessing)
		if err := s.repo.UpsertAsset(ctx, tx, asset); err != nil {
			utils.SetSpanError(ctx, err)
			return derrors.Infra("failed to update asset status", err)
		}
	}

	if len(jobMsg.Assets) == 0 {
		return derrors.Validation("no assets ready for processing", map[string]any{"listingIdentityId": input.ListingIdentityID})
	}

	// Send to Queue (which triggers Step Function)
	if _, err := s.queue.EnqueueJob(ctx, jobMsg); err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to publish job", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to commit process request", err)
	}
	committed = true

	logger.Info("service.media.process.started", "listing_identity_id", input.ListingIdentityID, "assets_count", len(assets), "job_id", jobID)
	return nil
}

func (s *mediaProcessingService) ensureRawObjectsExist(ctx context.Context, assets []mediaprocessingmodel.MediaAsset) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	seen := make(map[string]struct{})

	for _, asset := range assets {
		rawKey := strings.TrimSpace(asset.S3KeyRaw())
		if rawKey == "" {
			continue
		}
		if _, ok := seen[rawKey]; ok {
			continue
		}
		seen[rawKey] = struct{}{}

		if _, err := s.storage.ValidateObjectChecksum(ctx, rawKey, ""); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.process.raw_head_failed", "key", rawKey, "err", err)
			return derrors.Validation("raw media file not found or inaccessible", map[string]any{"key": rawKey})
		}
	}

	return nil
}
