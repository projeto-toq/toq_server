package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/media_processing_repository"
	"github.com/projeto-toq/toq_server/internal/core/port/right/workflow"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CompleteMedia finalizes the media processing, triggers ZIP generation and updates listing status.
func (s *mediaProcessingService) CompleteMedia(ctx context.Context, input dto.CompleteMediaInput) error {
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
		logger.Error("service.media.complete.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.complete.rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	// 1. Validate Listing Status
	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to load listing", err)
	}

	if listing.Status() != listingmodel.StatusPendingPhotoProcessing {
		return derrors.Conflict("listing is not in PENDING_PHOTO_PROCESSING status")
	}

	// 2. Fetch ALL assets to validate state
	// We do not filter by status initially to detect pending/processing items
	allAssets, err := s.repo.ListAssets(ctx, tx, uint64(input.ListingIdentityID), mediaprocessingrepository.AssetFilter{}, nil)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.list_assets_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to list assets", err)
	}

	processedAssets, err := ensureAssetsReadyForFinalization(allAssets)
	if err != nil {
		return err
	}

	// 4. Register ZIP Job
	job := mediaprocessingmodel.NewMediaProcessingJob(uint64(input.ListingIdentityID), mediaprocessingmodel.MediaProcessingProviderStepFunctions)
	jobID, err := s.repo.RegisterProcessingJob(ctx, tx, job)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.register_job_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to register zip job", err)
	}
	job.SetID(jobID)

	// 5. Trigger Finalization Pipeline
	finalizationInput := buildFinalizationInput(ctx, jobID, uint64(input.ListingIdentityID), processedAssets)

	executionARN, err := s.workflow.StartMediaFinalization(ctx, finalizationInput)
	if err != nil {
		if errors.Is(err, workflow.ErrFinalizationAccessDenied) {
			logger.Error("service.media.complete.workflow_denied", "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
			return derrors.Forbidden(
				"media finalization temporarily unavailable",
				derrors.WithDetails(map[string]any{"reason": "workflow_access_denied"}),
			)
		}

		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.start_workflow_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
		return derrors.Infra("failed to start finalization workflow", err)
	}

	// Update job with execution ARN
	job.MarkRunning(executionARN, s.now())
	if err := s.repo.UpdateProcessingJob(ctx, tx, job); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.update_job_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
		return derrors.Infra("failed to persist job state", err)
	}

	// 5. Advance Listing Status
	listing.SetStatus(listingmodel.StatusPendingOwnerApproval)
	if err := s.listingRepo.UpdateListingVersion(ctx, tx, listing); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.update_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return derrors.Infra("failed to update listing status", err)
	}

	// Commit
	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.commit_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	if err := s.notifyOwnerMediaReady(ctx, listing); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.complete.owner_notification_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", jobID)
	}

	logger.Info("service.media.complete.started_zip", "listing_identity_id", input.ListingIdentityID, "job_id", jobID, "execution_arn", executionARN, "assets_count", len(processedAssets))
	return nil
}
