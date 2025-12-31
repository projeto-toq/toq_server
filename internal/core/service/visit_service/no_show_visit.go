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

	visit, entry, hasEntry, err := s.loadVisitAndEntry(ctx, tx, visitID)
	if err != nil {
		return nil, err
	}
	if visit.Status() != listingmodel.VisitStatusApproved && visit.Status() != listingmodel.VisitStatusPending {
		return nil, utils.ConflictError("Only pending or approved visits can be marked as no-show")
	}

	visit.SetStatus(listingmodel.VisitStatusNoShow)
	if ownerNotes != "" {
		visit.SetOwnerNotes(ownerNotes)
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

	if hasEntry {
		entry.SetEntryType(schedulemodel.EntryTypeVisitConfirmed)
		entry.SetBlocking(false)
		if err = s.scheduleRepo.UpdateEntry(ctx, tx, entry); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("visit.no_show.update_entry_error", "entry_id", entry.ID(), "visit_id", visitID, "err", err)
			return nil, utils.InternalError("")
		}
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("visit.no_show.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}
	committed = true

	return visit, nil
}
