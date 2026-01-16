package scheduleservices

import (
	"context"
	"database/sql"
	"errors"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) FinishListingAgenda(ctx context.Context, input FinishListingAgendaInput) (err error) {
	if input.ListingIdentityID <= 0 {
		return utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if input.ActorID <= 0 {
		return utils.ValidationError("actorId", "actorId must be greater than zero")
	}

	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.finish_listing_agenda.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	defer func() {
		if err != nil {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.finish_listing_agenda.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	versionSummaries, err := s.listingRepo.ListListingVersions(ctx, tx, listingrepository.ListListingVersionsFilter{
		ListingIdentityID: input.ListingIdentityID,
		IncludeDeleted:    false,
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.finish_listing_agenda.list_versions_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	var activeVersion listingmodel.ListingVersionInterface
	for _, summary := range versionSummaries {
		if !summary.IsActive {
			continue
		}
		if summary.Version == nil {
			continue
		}
		activeVersion = summary.Version
		break
	}

	if activeVersion == nil {
		return utils.NotFoundError("Listing active version")
	}

	if activeVersion.UserID() != input.OwnerID {
		return utils.AuthorizationError("Only listing owner can finish agenda setup")
	}

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ConflictError("Listing agenda must be created before finishing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.finish_listing_agenda.get_agenda_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	if agenda == nil {
		return utils.ConflictError("Listing agenda must be created before finishing")
	}

	if agenda.OwnerID() != input.OwnerID {
		return utils.AuthorizationError("Agenda owner does not match user")
	}

	if activeVersion.Status() != listingmodel.StatusPendingAvailability {
		return utils.ConflictError("Listing must be pending availability to finish agenda")
	}

	updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, activeVersion.ID(), listingmodel.StatusPendingPhotoScheduling, listingmodel.StatusPendingAvailability)
	if updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return utils.ConflictError("Listing status changed while finishing agenda")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("schedule.finish_listing_agenda.update_status_error", "err", updateErr, "listing_identity_id", input.ListingIdentityID, "listing_version_id", activeVersion.ID())
		return utils.InternalError("")
	}

	agendaID := int64(agenda.ID())

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		input.ActorID,
		auditmodel.AuditTarget{Type: auditmodel.TargetListingAgenda, ID: agendaID},
		auditmodel.OperationAgendaFinish,
		map[string]any{
			"listing_identity_id": input.ListingIdentityID,
			"listing_version_id":  activeVersion.ID(),
			"agenda_id":           agendaID,
			"status_from":         listingmodel.StatusPendingAvailability.String(),
			"status_to":           listingmodel.StatusPendingPhotoScheduling.String(),
			"actor_role":          string(permissionmodel.RoleSlugOwner),
		},
	)

	if auditErr := s.auditService.RecordChange(ctx, tx, auditRecord); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("schedule.finish_listing_agenda.audit_error", "err", auditErr, "listing_identity_id", input.ListingIdentityID, "agenda_id", agendaID)
		return auditErr
	}

	if err = s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.finish_listing_agenda.tx_commit_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	logger.Info("schedule.finish_listing_agenda.completed", "listing_identity_id", input.ListingIdentityID, "listing_version_id", activeVersion.ID(), "new_status", listingmodel.StatusPendingPhotoScheduling.String())

	return nil
}
