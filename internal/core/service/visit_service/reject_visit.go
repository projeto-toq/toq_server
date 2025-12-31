package visitservice

import (
	"context"
	"database/sql"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *visitService) RejectVisit(ctx context.Context, visitID int64, reason string) (listingmodel.VisitInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if reason == "" {
		return nil, utils.ValidationError("reason", "Rejection reason is required")
	}

	actorID, uidErr := s.globalService.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return nil, uidErr
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("visit.reject.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("visit.reject.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	visit, entry, hasEntry, err := s.loadVisitAndEntry(ctx, tx, visitID)
	if err != nil {
		return nil, err
	}
	if visit.Status() != listingmodel.VisitStatusPending {
		return nil, utils.ConflictError("Only pending visits can be rejected")
	}

	now := time.Now().UTC()
	markFirstOwnerActionIfEmpty(visit, now)
	visit.SetStatus(listingmodel.VisitStatusRejected)
	visit.SetRejectionReason(reason)
	visit.SetUpdatedBy(actorID)

	if err = s.visitRepo.UpdateVisit(ctx, tx, visit); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("Visit")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.reject.update_visit_error", "visit_id", visitID, "err", err)
		return nil, utils.InternalError("")
	}

	if hasEntry {
		if err = s.scheduleRepo.DeleteEntry(ctx, tx, entry.ID()); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("visit.reject.delete_entry_error", "entry_id", entry.ID(), "visit_id", visitID, "err", err)
			return nil, utils.InternalError("")
		}
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("visit.reject.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}
	committed = true

	return visit, nil
}
