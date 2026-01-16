package listingservices

import (
	"context"
	"database/sql"
	"errors"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) DiscardDraftVersion(ctx context.Context, input DiscardDraftVersionInput) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate required fields
	if input.ListingIdentityID == 0 {
		return utils.BadRequest("listingIdentityId is required")
	}
	if input.VersionID == 0 {
		return utils.BadRequest("versionId is required")
	}

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.discard.tx_start_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.discard.tx_rollback_error", "err", rbErr, "version_id", input.VersionID)
			}
		}
	}()

	// Get user ID early for ownership validation
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return uidErr
	}

	// Get listing identity to validate ownership BEFORE fetching version
	identity, err := ls.listingRepository.GetListingIdentityByID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.discard.get_identity_error", "err", err, "identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	// Validate ownership using identity
	if identity.UserID != userID {
		logger.Warn("unauthorized_discard_attempt",
			"listing_identity_id", input.ListingIdentityID,
			"listing_version_id", input.VersionID,
			"requester_user_id", userID,
			"owner_user_id", identity.UserID)
		return utils.AuthorizationError("Only listing owner can discard a draft version")
	}

	// Now get the specific version
	draft, err := ls.listingRepository.GetListingVersionByID(ctx, tx, input.VersionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing version")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.discard.get_version_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("")
	}

	// Verify version belongs to the identity
	if draft.IdentityID() != input.ListingIdentityID {
		logger.Warn("version_identity_mismatch_discard",
			"listing_identity_id", input.ListingIdentityID,
			"listing_version_id", input.VersionID,
			"version_actual_identity_id", draft.IdentityID(),
			"requester_user_id", userID)
		return utils.BadRequest("version does not belong to specified listing")
	}

	if draft.Status() != listingmodel.StatusDraft {
		return utils.ConflictError("Only draft versions can be discarded")
	}

	if draft.ActiveVersionID() == draft.ID() {
		return utils.ConflictError("Cannot discard the active listing version")
	}

	draft.SetDeleted(true)

	if updateErr := ls.listingRepository.UpdateListingVersion(ctx, tx, draft); updateErr != nil && !errors.Is(updateErr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.discard.update_error", "err", updateErr, "version_id", draft.ID())
		return utils.InternalError("")
	}

	version := int64(draft.Version())
	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		userID,
		auditmodel.AuditTarget{Type: auditmodel.TargetListingIdentity, ID: draft.IdentityID(), Version: &version},
		auditmodel.OperationDiscard,
		map[string]any{
			"listing_identity_id": draft.IdentityID(),
			"listing_version_id":  draft.ID(),
			"version":             draft.Version(),
			"status_from":         draft.Status().String(),
			"status_to":           draft.Status().String(),
			"actor_role":          string(permissionmodel.RoleSlugOwner),
			"deleted":             draft.Deleted(),
		},
	)

	if auditErr := ls.auditService.RecordChange(ctx, tx, auditRecord); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("listing.discard.audit_error", "err", auditErr, "listing_identity_id", draft.IdentityID(), "listing_version_id", draft.ID())
		return auditErr
	}

	if err = ls.gsi.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.discard.tx_commit_error", "err", err, "version_id", draft.ID())
		return utils.InternalError("")
	}

	logger.Info("listing.discard.completed", "version_id", draft.ID(), "identity_id", draft.IdentityID())

	return nil
}
