package holidayservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *holidayService) UpdateCalendar(ctx context.Context, input UpdateCalendarInput) (holidaymodel.CalendarInterface, error) {
	if input.ID == 0 {
		return nil, utils.ValidationError("id", "id must be greater than zero")
	}
	if strings.TrimSpace(input.Name) == "" {
		return nil, utils.ValidationError("name", "name is required")
	}
	if err := validateScopeInput(input.Scope, input.State, input.CityIBGE); err != nil {
		return nil, err
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
		logger.Error("holiday.update_calendar.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("holiday.update_calendar.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	existing, err := s.repo.GetCalendarByID(ctx, tx, input.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Holiday calendar")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.update_calendar.get_error", "id", input.ID, "err", err)
		return nil, utils.InternalError("")
	}

	existing.SetName(strings.TrimSpace(input.Name))
	existing.SetScope(input.Scope)
	if state := cleanState(input.Scope, input.State); state != "" {
		existing.SetState(state)
	} else {
		existing.ClearState()
	}
	if city := strings.TrimSpace(input.CityIBGE); city != "" {
		existing.SetCityIBGE(city)
	} else {
		existing.ClearCityIBGE()
	}
	existing.SetActive(input.IsActive)
	if err := s.repo.UpdateCalendar(ctx, tx, existing); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.update_calendar.repo_error", "id", input.ID, "err", err)
		return nil, utils.InternalError("")
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("holiday.update_calendar.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	committed = true
	return existing, nil
}
