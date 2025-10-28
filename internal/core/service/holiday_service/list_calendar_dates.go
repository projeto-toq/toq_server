package holidayservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *holidayService) ListCalendarDates(ctx context.Context, filter holidaymodel.CalendarDatesFilter) (holidaymodel.CalendarDatesResult, error) {
	if filter.CalendarID == 0 {
		return holidaymodel.CalendarDatesResult{}, utils.ValidationError("calendarId", "calendarId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return holidaymodel.CalendarDatesResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("holiday.list_calendar_dates.tx_start_error", "err", txErr)
		return holidaymodel.CalendarDatesResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("holiday.list_calendar_dates.tx_rollback_error", "err", rbErr)
		}
	}()

	calendar, err := s.repo.GetCalendarByID(ctx, tx, filter.CalendarID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return holidaymodel.CalendarDatesResult{}, utils.NotFoundError("Holiday calendar")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.list_calendar_dates.get_calendar_error", "calendar_id", filter.CalendarID, "err", err)
		return holidaymodel.CalendarDatesResult{}, utils.InternalError("")
	}

	tzName := filter.Timezone
	if strings.TrimSpace(tzName) == "" {
		tzName = calendar.Timezone()
	}
	loc, tzErr := utils.ResolveLocation("timezone", tzName)
	if tzErr != nil {
		return holidaymodel.CalendarDatesResult{}, tzErr
	}

	repoFilter := filter
	repoFilter.From, repoFilter.To = utils.NormalizeRangeToUTC(filter.From, filter.To, loc)

	result, err := s.repo.ListCalendarDates(ctx, tx, repoFilter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.list_calendar_dates.repo_error", "calendar_id", filter.CalendarID, "err", err)
		return holidaymodel.CalendarDatesResult{}, utils.InternalError("")
	}

	for _, date := range result.Dates {
		date.SetHolidayDate(utils.ConvertToLocation(date.HolidayDate(), loc))
		date.SetTimezone(loc.String())
	}

	return result, nil
}
