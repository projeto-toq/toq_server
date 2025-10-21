package holidayservices

import (
	"context"
	"database/sql"
	"errors"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *holidayService) GetCalendarByID(ctx context.Context, id uint64) (holidaymodel.CalendarInterface, error) {
	if id == 0 {
		return nil, utils.ValidationError("id", "id must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("holiday.get_calendar.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("holiday.get_calendar.tx_rollback_error", "err", rbErr)
		}
	}()

	calendar, err := s.repo.GetCalendarByID(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Holiday calendar")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.get_calendar.repo_error", "id", id, "err", err)
		return nil, utils.InternalError("")
	}

	return calendar, nil
}
