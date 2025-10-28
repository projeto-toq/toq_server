package holidayservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *holidayService) CreateCalendarDate(ctx context.Context, input CreateCalendarDateInput) (holidaymodel.CalendarDateInterface, error) {
	if input.CalendarID == 0 {
		return nil, utils.ValidationError("calendarId", "calendarId must be greater than zero")
	}
	if input.HolidayDate.IsZero() {
		return nil, utils.ValidationError("holidayDate", "holidayDate is required")
	}
	if strings.TrimSpace(input.Label) == "" {
		return nil, utils.ValidationError("label", "label is required")
	}
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("holiday.create_calendar_date.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("holiday.create_calendar_date.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	calendar, err := s.repo.GetCalendarByID(ctx, tx, input.CalendarID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Holiday calendar")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.create_calendar_date.get_calendar_error", "calendar_id", input.CalendarID, "err", err)
		return nil, utils.InternalError("")
	}

	loc, tzErr := utils.ResolveLocation("timezone", calendar.Timezone())
	if tzErr != nil {
		return nil, tzErr
	}

	domain := holidaymodel.NewCalendarDate()
	domain.SetCalendarID(input.CalendarID)
	holidayAtLoc := utils.NormalizeDateToLocationMidnight(input.HolidayDate, loc)
	domain.SetHolidayDate(holidayAtLoc)
	domain.SetTimezone(loc.String())
	domain.SetLabel(strings.TrimSpace(input.Label))
	domain.SetRecurrent(input.Recurrent)

	id, err := s.repo.CreateCalendarDate(ctx, tx, domain)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.create_calendar_date.repo_error", "calendar_id", input.CalendarID, "err", err)
		return nil, utils.InternalError("")
	}
	domain.SetID(id)
	domain.SetHolidayDate(utils.ConvertToLocation(domain.HolidayDate(), loc))

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("holiday.create_calendar_date.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	committed = true
	return domain, nil
}
