package visitservice

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *visitService) CancelVisit(ctx context.Context, visitID int64, reason string) (listingmodel.VisitInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if reason == "" {
		return nil, utils.ValidationError("reason", "Cancel reason is required")
	}

	actorID, uidErr := s.globalService.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return nil, uidErr
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("visit.cancel.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("visit.cancel.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	visit, entry, hasEntry, err := s.loadVisitAndEntry(ctx, tx, visitID)
	if err != nil {
		return nil, err
	}
	if visit.Status() != listingmodel.VisitStatusPending && visit.Status() != listingmodel.VisitStatusApproved {
		return nil, utils.ConflictError("Only pending or approved visits can be cancelled")
	}

	visit.SetStatus(listingmodel.VisitStatusCancelled)
	visit.SetCancelReason(reason)
	visit.SetUpdatedBy(actorID)

	if err = s.visitRepo.UpdateVisit(ctx, tx, visit); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("Visit")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.cancel.update_visit_error", "visit_id", visitID, "err", err)
		return nil, utils.InternalError("")
	}

	if hasEntry {
		if err = s.scheduleRepo.DeleteEntry(ctx, tx, entry.ID()); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("visit.cancel.delete_entry_error", "entry_id", entry.ID(), "visit_id", visitID, "err", err)
			return nil, utils.InternalError("")
		}
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("visit.cancel.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}
	committed = true

	return visit, nil
}
