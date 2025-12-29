package listingservices

import (
	"context"
	"database/sql"
	"errors"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ChangeListingStatus performs owner-driven status transitions between READY and PUBLISHED states.
//
// Flow:
//  1. Validate payload and ensure requester ownership
//  2. Fetch active version for the listing identity
//  3. Validate transition rules and update status atomically
//  4. Register audit trail and commit transaction
func (ls *listingService) ChangeListingStatus(ctx context.Context, input ChangeListingStatusInput) (output ChangeListingStatusOutput, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return output, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err = input.Validate(); err != nil {
		return output, err
	}

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.change_status.tx_start_error", "err", err)
		return output, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.change_status.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	identity, err := ls.listingRepository.GetListingIdentityByID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return output, utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.change_status.get_identity_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	if identity.Deleted {
		return output, utils.NotFoundError("listing")
	}

	if identity.UserID != input.RequesterUserID {
		logger.Warn("listing.change_status.unauthorized", "listing_identity_id", identity.ID, "owner_id", identity.UserID, "requester_id", input.RequesterUserID)
		return output, utils.AuthorizationError("Only listing owner can change status")
	}

	activeVersion, err := ls.listingRepository.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return output, utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.change_status.get_active_version_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	previousStatus := activeVersion.Status()
	targetStatus, validationErr := resolveTargetStatus(previousStatus, input.Action)
	if validationErr != nil {
		return output, validationErr
	}

	if err = ls.listingRepository.UpdateListingStatus(ctx, tx, activeVersion.ID(), targetStatus, previousStatus); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return output, utils.ConflictError("Listing status changed while processing request")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.change_status.update_status_error", "err", err, "listing_version_id", activeVersion.ID(), "target_status", targetStatus.String())
		return output, utils.InternalError("")
	}

	if auditErr := ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Listing status changed by owner"); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("listing.change_status.audit_error", "err", auditErr, "listing_identity_id", input.ListingIdentityID)
		return output, auditErr
	}

	if err = ls.gsi.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.change_status.tx_commit_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	output = ChangeListingStatusOutput{
		ListingIdentityID: input.ListingIdentityID,
		ActiveVersionID:   activeVersion.ID(),
		PreviousStatus:    previousStatus,
		NewStatus:         targetStatus,
	}

	logger.Info("listing.change_status.completed",
		"listing_identity_id", input.ListingIdentityID,
		"listing_version_id", activeVersion.ID(),
		"previous_status", previousStatus.String(),
		"new_status", targetStatus.String(),
		"action", string(input.Action))

	return output, nil
}

// resolveTargetStatus validates the current status versus the desired action and returns the result.
func resolveTargetStatus(current listingmodel.ListingStatus, action ListingStatusAction) (listingmodel.ListingStatus, error) {
	switch action {
	case ListingStatusActionPublish:
		if current != listingmodel.StatusReady {
			return 0, utils.BadRequest("Listing must be READY to publish")
		}
		return listingmodel.StatusPublished, nil
	case ListingStatusActionSuspend:
		if current != listingmodel.StatusPublished && current != listingmodel.StatusUnderOffer && current != listingmodel.StatusUnderNegotiation {
			return 0, utils.BadRequest("Listing must be PUBLISHED, UNDER_OFFER or UNDER_NEGOTIATION to suspend")
		}
		return listingmodel.StatusReady, nil
	default:
		return 0, utils.ValidationError("action", "Unsupported action")
	}
}
