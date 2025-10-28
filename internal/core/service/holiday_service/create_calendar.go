package holidayservices

import (
	"context"
	"strings"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *holidayService) CreateCalendar(ctx context.Context, input CreateCalendarInput) (holidaymodel.CalendarInterface, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, utils.ValidationError("name", "name is required")
	}
	if err := validateScopeInput(input.Scope, input.State, input.City); err != nil {
		return nil, err
	}
	loc, tzErr := utils.ResolveLocation("timezone", input.Timezone)
	if tzErr != nil {
		return nil, tzErr
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
		logger.Error("holiday.create_calendar.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("holiday.create_calendar.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	domain := holidaymodel.NewCalendar()
	domain.SetName(strings.TrimSpace(input.Name))
	domain.SetScope(input.Scope)
	if state := cleanState(input.Scope, input.State); state != "" {
		domain.SetState(state)
	}
	if city := strings.TrimSpace(input.City); city != "" {
		domain.SetCity(city)
	}
	domain.SetActive(input.IsActive || input.Scope == holidaymodel.ScopeNational)
	domain.SetTimezone(loc.String())

	id, err := s.repo.CreateCalendar(ctx, tx, domain)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.create_calendar.repo_error", "err", err)
		return nil, utils.InternalError("")
	}
	domain.SetID(id)

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("holiday.create_calendar.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	committed = true
	return domain, nil
}
