package holidayservices

import (
	"context"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *holidayService) ListCalendars(ctx context.Context, filter holidaymodel.CalendarListFilter) (holidaymodel.CalendarListResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return holidaymodel.CalendarListResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("holiday.list_calendars.tx_start_error", "err", txErr)
		return holidaymodel.CalendarListResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("holiday.list_calendars.tx_rollback_error", "err", rbErr)
		}
	}()

	result, err := s.repo.ListCalendars(ctx, tx, filter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.list_calendars.repo_error", "err", err)
		return holidaymodel.CalendarListResult{}, utils.InternalError("")
	}

	return result, nil
}
