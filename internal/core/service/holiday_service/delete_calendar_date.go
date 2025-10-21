package holidayservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *holidayService) DeleteCalendarDate(ctx context.Context, id uint64) error {
	if id == 0 {
		return utils.ValidationError("id", "id must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("holiday.delete_calendar_date.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("holiday.delete_calendar_date.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if err := s.repo.DeleteCalendarDate(ctx, tx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Holiday calendar date")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.delete_calendar_date.repo_error", "id", id, "err", err)
		return utils.InternalError("")
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("holiday.delete_calendar_date.tx_commit_error", "err", cmErr)
		return utils.InternalError("")
	}

	committed = true
	return nil
}
