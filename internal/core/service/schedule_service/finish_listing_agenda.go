package scheduleservices

import (
	"context"
	"database/sql"
	"errors"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) FinishListingAgenda(ctx context.Context, input FinishListingAgendaInput) (err error) {
	if input.ListingID <= 0 {
		return utils.ValidationError("listingId", "listingId must be greater than zero")
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
		logger.Error("schedule.finish_listing_agenda.tx_start_error", "err", txErr, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	defer func() {
		if err != nil {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.finish_listing_agenda.tx_rollback_error", "err", rbErr, "listing_id", input.ListingID)
			}
		}
	}()

	listing, err := s.listingRepo.GetListingByID(ctx, tx, input.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.finish_listing_agenda.get_listing_error", "err", err, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	if listing.UserID() != input.OwnerID {
		return utils.AuthorizationError("Only listing owner can finish agenda setup")
	}

	agenda, err := s.scheduleRepo.GetAgendaByListingID(ctx, tx, input.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ConflictError("Listing agenda must be created before finishing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.finish_listing_agenda.get_agenda_error", "err", err, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	if agenda == nil {
		return utils.ConflictError("Listing agenda must be created before finishing")
	}

	if agenda.OwnerID() != input.OwnerID {
		return utils.AuthorizationError("Agenda owner does not match user")
	}

	if listing.Status() != listingmodel.StatusPendingAvailability {
		return utils.ConflictError("Listing must be pending availability to finish agenda")
	}

	updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, input.ListingID, listingmodel.StatusPendingPhotoScheduling, listingmodel.StatusPendingAvailability)
	if updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return utils.ConflictError("Listing status changed while finishing agenda")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("schedule.finish_listing_agenda.update_status_error", "err", updateErr, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	if auditErr := s.globalService.CreateAudit(ctx, tx, globalmodel.TableListings, "Agenda finalizada (schedule.finish)", input.ActorID); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("schedule.finish_listing_agenda.audit_error", "err", auditErr, "listing_id", input.ListingID)
		return auditErr
	}

	if err = s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.finish_listing_agenda.tx_commit_error", "err", err, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	logger.Info("schedule.finish_listing_agenda.completed", "listing_id", input.ListingID, "new_status", listingmodel.StatusPendingPhotoScheduling.String())

	return nil
}
