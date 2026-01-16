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
	"github.com/projeto-toq/toq_server/internal/core/port/right/workflow"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CompleteProjectMedia finalizes project media by copying raw assets to processed paths, registering a ZIP job,
// and advancing the listing status to admin review or ready.
func (s *mediaProcessingService) CompleteProjectMedia(ctx context.Context, input dto.CompleteMediaInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if !s.cfg.AllowOwnerProjectUpload {
		return derrors.Forbidden("project uploads are disabled", nil)
	}

	if input.ListingIdentityID == 0 {
		return derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.project_complete.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.project_complete.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_complete.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to load listing", err)
	}

	if err := s.ensureProjectFlowAllowed(listing); err != nil {
		return err
	}

	filter := mediaprocessingrepository.AssetFilter{
		AssetTypes: []mediaprocessingmodel.MediaAssetType{
			mediaprocessingmodel.MediaAssetTypeProjectDoc,
			mediaprocessingmodel.MediaAssetTypeProjectRender,
		},
	}

	assets, err := s.repo.ListAssets(ctx, tx, uint64(input.ListingIdentityID), filter, nil)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_complete.list_assets_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to list project assets", err)
	}
	if len(assets) == 0 {
		return derrors.Conflict("no project assets found for this listing")
	}

	processedAssets := make([]mediaprocessingmodel.MediaAsset, 0, len(assets))
	for idx := range assets {
		asset := assets[idx]
		status := asset.Status()

		rawKey := strings.TrimSpace(asset.S3KeyRaw())
		if rawKey == "" {
			return derrors.Conflict("project asset is missing raw key", derrors.WithDetails(map[string]any{"assetType": asset.AssetType(), "sequence": asset.Sequence()}))
		}

		switch status {
		case mediaprocessingmodel.MediaAssetStatusPendingUpload:
			return derrors.Conflict("project asset is still pending upload", derrors.WithDetails(map[string]any{"assetType": asset.AssetType(), "sequence": asset.Sequence()}))
		case mediaprocessingmodel.MediaAssetStatusFailed:
			return derrors.Conflict("project asset processing failed; re-upload required", derrors.WithDetails(map[string]any{"assetType": asset.AssetType(), "sequence": asset.Sequence()}))
		case mediaprocessingmodel.MediaAssetStatusProcessing:
			return derrors.Conflict("project asset is processing; please wait", derrors.WithDetails(map[string]any{"assetType": asset.AssetType(), "sequence": asset.Sequence()}))
		case mediaprocessingmodel.MediaAssetStatusProcessed:
			if strings.TrimSpace(asset.S3KeyProcessed()) == "" {
				return derrors.Conflict("processed project asset missing processed key", derrors.WithDetails(map[string]any{"assetType": asset.AssetType(), "sequence": asset.Sequence()}))
			}
			processedAssets = append(processedAssets, asset)
			continue
		}

		processedKey, copyErr := s.copyProjectAssetToProcessed(ctx, uint64(input.ListingIdentityID), asset)
		if copyErr != nil {
			return copyErr
		}

		asset.SetS3KeyProcessed(processedKey)
		asset.SetStatus(mediaprocessingmodel.MediaAssetStatusProcessed)

		if err := s.repo.UpsertAsset(ctx, tx, asset); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.project_complete.upsert_asset_error", "err", err, "listing_identity_id", input.ListingIdentityID)
			return derrors.Infra("failed to persist processed asset", err)
		}

		processedAssets = append(processedAssets, asset)
	}

	if len(processedAssets) == 0 {
		return derrors.Conflict("no processed project assets available")
	}

	job := mediaprocessingmodel.NewMediaProcessingJob(uint64(input.ListingIdentityID), mediaprocessingmodel.MediaProcessingProviderStepFunctionsFinalization)
	jobID, err := s.repo.RegisterProcessingJob(ctx, tx, job)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_complete.register_job_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to register finalization job", err)
	}
	job.SetID(jobID)

	finalizationInput := buildFinalizationInput(ctx, jobID, uint64(input.ListingIdentityID), processedAssets)
	executionARN, err := s.workflow.StartMediaFinalization(ctx, finalizationInput)
	if err != nil {
		if errors.Is(err, workflow.ErrFinalizationAccessDenied) {
			logger.Error("service.media.project_complete.workflow_denied", "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
			return derrors.Forbidden("media finalization temporarily unavailable", derrors.WithDetails(map[string]any{"reason": "workflow_access_denied"}))
		}

		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_complete.start_workflow_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
		return derrors.Infra("failed to start finalization workflow", err)
	}

	job.MarkRunning(executionARN, s.now())
	if err := s.repo.UpdateProcessingJob(ctx, tx, job); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_complete.update_job_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
		return derrors.Infra("failed to persist job state", err)
	}

	if s.cfg.RequireAdminReview {
		listing.SetStatus(listingmodel.StatusPendingAdminReview)
	} else {
		listing.SetStatus(listingmodel.StatusReady)
	}

	if err := s.listingRepo.UpdateListingVersion(ctx, tx, listing); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_complete.update_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to update listing status", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_complete.commit_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	logger.Info("service.media.project_complete.started_zip", "listing_identity_id", input.ListingIdentityID, "job_id", jobID, "assets_count", len(processedAssets))
	return nil
}
