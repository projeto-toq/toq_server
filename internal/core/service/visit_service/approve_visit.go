package visitservice

import (
	"context"
	"database/sql"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *visitService) ApproveVisit(ctx context.Context, visitID int64, ownerNotes string) (listingmodel.VisitInterface, error) {
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
		logger.Error("visit.approve.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("visit.approve.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	visit, err := s.loadVisit(ctx, tx, visitID)
	if err != nil {
		return nil, err
	}
	if visit.Status() != listingmodel.VisitStatusPending {
		return nil, utils.ConflictError("Only pending visits can be approved")
	}

	agenda, agErr := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, visit.ListingIdentityID())
	if agErr != nil {
		if agErr == sql.ErrNoRows {
			return nil, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, agErr)
		logger.Error("visit.approve.get_agenda_error", "listing_identity_id", visit.ListingIdentityID(), "err", agErr)
		return nil, utils.InternalError("")
	}

	now := time.Now().UTC()
	if err := s.recordOwnerResponseMetrics(ctx, tx, visit, actorID, now); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.approve.owner_response_metrics_error", "visit_id", visitID, "err", err)
		return nil, err
	}
	visit.SetStatus(listingmodel.VisitStatusApproved)
	if ownerNotes != "" {
		visit.SetOwnerNotes(ownerNotes)
	}
	visit.SetUpdatedBy(actorID)

	if err = s.visitRepo.UpdateVisit(ctx, tx, visit); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("Visit")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.approve.update_visit_error", "visit_id", visitID, "err", err)
		return nil, utils.InternalError("")
	}

	if err = s.ensureVisitEntries(ctx, tx, agenda, visit, schedulemodel.EntryTypeVisitConfirmed, true); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.approve.ensure_entries_error", "visit_id", visitID, "err", err)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("visit.approve.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}
	committed = true

	s.notifyVisitStatusOwner(ctx, visit)
	s.notifyVisitStatusRealtor(ctx, visit)

	return visit, nil
}
