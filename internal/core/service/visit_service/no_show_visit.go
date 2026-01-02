package visitservice

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *visitService) MarkNoShow(ctx context.Context, visitID int64, ownerNotes string) (listingmodel.VisitInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	actorID, uidErr := s.globalService.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return nil, uidErr
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("visit.no_show.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("visit.no_show.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	visit, err := s.loadVisit(ctx, tx, visitID)
	if err != nil {
		return nil, err
	}
	if visit.Status() != listingmodel.VisitStatusApproved && visit.Status() != listingmodel.VisitStatusPending {
		return nil, utils.ConflictError("Only pending or approved visits can be marked as no-show")
	}

	agenda, agErr := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, visit.ListingIdentityID())
	if agErr != nil {
		if agErr == sql.ErrNoRows {
			return nil, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, agErr)
		logger.Error("visit.no_show.get_agenda_error", "listing_identity_id", visit.ListingIdentityID(), "err", agErr)
		return nil, utils.InternalError("")
	}

	visit.SetStatus(listingmodel.VisitStatusNoShow)
	if ownerNotes != "" {
		visit.SetNotes(ownerNotes)
	}
	visit.SetUpdatedBy(actorID)

	if err = s.visitRepo.UpdateVisit(ctx, tx, visit); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("Visit")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.no_show.update_visit_error", "visit_id", visitID, "err", err)
		return nil, utils.InternalError("")
	}

	if err = s.ensureVisitEntries(ctx, tx, agenda, visit, schedulemodel.EntryTypeVisitConfirmed, false); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.no_show.ensure_entries_error", "visit_id", visitID, "err", err)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("visit.no_show.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}
	committed = true

	s.notifyVisitStatusOwner(ctx, visit)
	s.notifyVisitStatusRealtor(ctx, visit)

	return visit, nil
}
