package listingservices

import (
	"context"
	"database/sql"
	"errors"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) PromoteListingVersion(ctx context.Context, input PromoteListingVersionInput) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate required fields
	if input.ListingIdentityID == 0 {
		return utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if input.VersionID <= 0 {
		return utils.ValidationError("versionId", "versionId must be greater than zero")
	}

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.promote.tx_start_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.promote.tx_rollback_error", "err", rbErr, "version_id", input.VersionID)
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
		logger.Error("listing.promote.get_identity_error", "err", err, "identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	// Validate ownership using identity
	if identity.UserID != userID {
		logger.Warn("unauthorized_promote_attempt",
			"listing_identity_id", input.ListingIdentityID,
			"listing_version_id", input.VersionID,
			"requester_user_id", userID,
			"owner_user_id", identity.UserID)
		return utils.AuthorizationError("Only listing owner can promote version")
	}

	listingVersionID := input.VersionID

	snapshot, err := ls.listingRepository.GetListingForEndUpdate(ctx, tx, listingVersionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.promote.fetch_error", "err", err, "listing_version_id", listingVersionID)
		return utils.InternalError("")
	}

	// Verify version belongs to the identity
	if snapshot.ListingID != input.ListingIdentityID {
		logger.Warn("version_identity_mismatch_promote",
			"listing_identity_id", input.ListingIdentityID,
			"listing_version_id", input.VersionID,
			"version_actual_identity_id", snapshot.ListingID,
			"requester_user_id", userID)
		return utils.BadRequest("version does not belong to specified listing")
	}

	if snapshot.Status != listingmodel.StatusDraft {
		return utils.ConflictError("Only draft versions can be promoted")
	}

	if verr := ls.validateListingBeforeEndUpdate(ctx, tx, snapshot); verr != nil {
		return verr
	}

	// Determine target status based on version number
	var targetStatus listingmodel.ListingStatus
	if snapshot.Version == 1 {
		// V1 promotion: Always start at PendingAvailability
		targetStatus = listingmodel.StatusPendingAvailability
	} else {
		// V>1 promotion: Inherit status from previous active version
		previousStatus, statusErr := ls.listingRepository.GetPreviousActiveVersionStatus(ctx, tx, snapshot.ListingID)
		if statusErr != nil {
			if errors.Is(statusErr, sql.ErrNoRows) {
				// Fallback: no previous active version found, use PendingAvailability
				logger.Warn("listing.promote.no_previous_active", "listing_identity_id", snapshot.ListingID, "version", snapshot.Version)
				targetStatus = listingmodel.StatusPendingAvailability
			} else {
				utils.SetSpanError(ctx, statusErr)
				logger.Error("listing.promote.get_previous_status_error", "err", statusErr, "listing_identity_id", snapshot.ListingID)
				return utils.InternalError("")
			}
		} else {
			targetStatus = previousStatus
		}
	}

	updateErr := ls.listingRepository.UpdateListingStatus(ctx, tx, listingVersionID, targetStatus, listingmodel.StatusDraft)
	if updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return utils.ConflictError("Listing status changed while promoting version")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.promote.update_status_error", "err", updateErr, "listing_version_id", listingVersionID, "target_status", targetStatus)
		return utils.InternalError("")
	}

	// Set the new active version
	if setErr := ls.listingRepository.SetListingActiveVersion(ctx, tx, snapshot.ListingID, listingVersionID); setErr != nil {
		utils.SetSpanError(ctx, setErr)
		logger.Error("listing.promote.set_active_error", "err", setErr, "listing_identity_id", snapshot.ListingID, "version_id", listingVersionID)
		return utils.InternalError("")
	}

	version := int64(snapshot.Version)
	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		userID,
		auditmodel.AuditTarget{Type: auditmodel.TargetListingIdentity, ID: snapshot.ListingID, Version: &version},
		auditmodel.OperationPromote,
		map[string]any{
			"listing_identity_id": snapshot.ListingID,
			"listing_version_id":  listingVersionID,
			"version":             snapshot.Version,
			"status_from":         listingmodel.StatusDraft.String(),
			"status_to":           targetStatus.String(),
			"actor_role":          string(permissionmodel.RoleSlugOwner),
		},
	)

	if auditErr := ls.auditService.RecordChange(ctx, tx, auditRecord); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("listing.promote.audit_error", "err", auditErr, "listing_identity_id", snapshot.ListingID, "listing_version_id", listingVersionID)
		return auditErr
	}

	// If V1 or target status is PendingAvailability, create default agenda
	if snapshot.Version == 1 || targetStatus == listingmodel.StatusPendingAvailability {
		timezone := resolveListingTimezone(snapshot)
		agendaInput := scheduleservices.CreateDefaultAgendaInput{
			ListingIdentityID: snapshot.ListingID,
			OwnerID:           userID,
			Timezone:          timezone,
			ActorID:           userID,
		}
		if _, agendaErr := ls.scheduleService.CreateDefaultAgendaWithTx(ctx, tx, agendaInput); agendaErr != nil {
			utils.SetSpanError(ctx, agendaErr)
			logger.Error("listing.promote.create_default_agenda_error", "err", agendaErr, "listing_identity_id", snapshot.ListingID)
			return agendaErr
		}
	}

	if err = ls.gsi.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.promote.tx_commit_error", "err", err, "listing_version_id", listingVersionID)
		return utils.InternalError("")
	}

	logger.Info("listing.promote.completed", "listing_version_id", listingVersionID, "listing_identity_id", snapshot.ListingID, "new_status", targetStatus.String())

	return nil
}
