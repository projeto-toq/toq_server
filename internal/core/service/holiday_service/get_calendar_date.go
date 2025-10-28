package holidayservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetCalendarDateByID returns a calendar date using the provided identifier and timezone preference.
func (s *holidayService) GetCalendarDateByID(ctx context.Context, id uint64, timezone string) (holidaymodel.CalendarDateInterface, error) {
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
		logger.Error("holiday.get_calendar_date.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("holiday.get_calendar_date.tx_rollback_error", "err", rbErr)
		}
	}()

	date, err := s.repo.GetCalendarDateByID(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Holiday date")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.get_calendar_date.repo_error", "date_id", id, "err", err)
		return nil, utils.InternalError("")
	}

	calendar, err := s.repo.GetCalendarByID(ctx, tx, date.CalendarID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Holiday calendar")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.get_calendar_date.get_calendar_error", "calendar_id", date.CalendarID(), "err", err)
		return nil, utils.InternalError("")
	}

	tzName := strings.TrimSpace(timezone)
	if tzName == "" {
		tzName = calendar.Timezone()
	}

	loc, tzErr := utils.ResolveLocation("timezone", tzName)
	if tzErr != nil {
		return nil, tzErr
	}

	date.SetHolidayDate(utils.ConvertToLocation(date.HolidayDate(), loc))
	date.SetTimezone(loc.String())

	return date, nil
}
