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
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CompleteMedia finalizes the media processing, triggers ZIP generation and updates listing status.
func (s *mediaProcessingService) CompleteMedia(ctx context.Context, input dto.CompleteMediaInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		return derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("service.media.complete.rollback_error", "err", rbErr)
			}
		}
	}()

	// 1. Validate Listing Status
	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("listing not found")
		}
		return derrors.Infra("failed to load listing", err)
	}

	if listing.Status() != listingmodel.StatusPendingPhotoProcessing {
		return derrors.Conflict("listing is not in PENDING_PHOTO_PROCESSING status")
	}

	// 2. Fetch processed assets for ZIP
	filter := mediaprocessingrepository.AssetFilter{
		Status: []mediaprocessingmodel.MediaAssetStatus{mediaprocessingmodel.MediaAssetStatusProcessed},
	}
	assets, err := s.repo.ListAssets(ctx, tx, uint64(input.ListingIdentityID), filter)
	if err != nil {
		return derrors.Infra("failed to list processed assets", err)
	}

	if len(assets) == 0 {
		return derrors.Conflict("no processed assets found to finalize")
	}

	// 3. Register ZIP Job
	job := mediaprocessingmodel.NewMediaProcessingJob(uint64(input.ListingIdentityID), mediaprocessingmodel.MediaProcessingProviderStepFunctions)
	jobID, err := s.repo.RegisterProcessingJob(ctx, tx, job)
	if err != nil {
		return derrors.Infra("failed to register zip job", err)
	}

	// 4. Trigger Finalization Pipeline
	jobAssets := make([]mediaprocessingmodel.JobAsset, 0, len(assets))
	for _, a := range assets {
		jobAssets = append(jobAssets, mediaprocessingmodel.JobAsset{
			Key:  a.S3KeyProcessed(),
			Type: string(a.AssetType()),
		})
	}

	finalizationInput := mediaprocessingmodel.MediaFinalizationInput{
		JobID:     jobID,
		ListingID: uint64(input.ListingIdentityID),
		Assets:    jobAssets,
	}

	executionARN, err := s.workflow.StartMediaFinalization(ctx, finalizationInput)
	if err != nil {
		return derrors.Infra("failed to start finalization workflow", err)
	}

	// Update job with execution ARN
	job.MarkRunning(executionARN, s.now())
	if err := s.repo.UpdateProcessingJob(ctx, tx, job); err != nil {
		logger.Error("service.media.complete.update_job_error", "err", err)
		// Non-critical, workflow is running
	}

	// 5. Advance Listing Status
	listing.SetStatus(listingmodel.StatusPendingOwnerApproval)
	if err := s.listingRepo.UpdateListingVersion(ctx, tx, listing); err != nil {
		return derrors.Infra("failed to update listing status", err)
	}

	// Commit
	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	return nil
}
