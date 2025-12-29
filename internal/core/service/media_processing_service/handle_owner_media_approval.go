package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HandleOwnerMediaApproval applies the owner's decision to the listing lifecycle.
func (s *mediaProcessingService) HandleOwnerMediaApproval(
	ctx context.Context,
	input dto.ListingMediaApprovalInput,
) (dto.ListingMediaApprovalOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return dto.ListingMediaApprovalOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID <= 0 {
		return dto.ListingMediaApprovalOutput{}, derrors.Validation(
			"listingIdentityId must be greater than zero",
			map[string]any{"listingIdentityId": "required"},
		)
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.owner_approval.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return dto.ListingMediaApprovalOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.owner_approval.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.ListingMediaApprovalOutput{}, derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.owner_approval.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return dto.ListingMediaApprovalOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Status() != listingmodel.StatusPendingOwnerApproval {
		return dto.ListingMediaApprovalOutput{}, derrors.BadRequest("listing is not awaiting owner approval", nil)
	}

	if listing.UserID() != int64(input.RequestedBy) {
		return dto.ListingMediaApprovalOutput{}, derrors.Forbidden("only the listing owner can approve or reject media")
	}

	var targetStatus listingmodel.ListingStatus
	if input.Approve {
		if s.cfg.RequireAdminReview {
			targetStatus = listingmodel.StatusPendingAdminReview
		} else {
			targetStatus = listingmodel.StatusReady
		}
	} else {
		targetStatus = listingmodel.StatusRejectedByOwner
	}

	listing.SetStatus(targetStatus)
	if err := s.listingRepo.UpdateListingVersion(ctx, tx, listing); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.owner_approval.update_status_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return dto.ListingMediaApprovalOutput{}, derrors.Infra("failed to update listing status", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.owner_approval.commit_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return dto.ListingMediaApprovalOutput{}, derrors.Infra("failed to commit owner approval", err)
	}
	committed = true

	logger.Info("service.media.owner_approval.applied",
		"listing_identity_id", input.ListingIdentityID,
		"decision", map[bool]string{true: "approved", false: "rejected"}[input.Approve],
		"new_status", targetStatus.String(),
	)

	return dto.ListingMediaApprovalOutput{
		ListingIdentityID: input.ListingIdentityID,
		NewStatus:         targetStatus,
	}, nil
}
