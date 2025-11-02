package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/converters"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *HolidayAdapter) CreateCalendar(ctx context.Context, tx *sql.Tx, calendar holidaymodel.CalendarInterface) (uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToCalendarEntity(calendar)

	query := `INSERT INTO holiday_calendars (name, scope, state, city, is_active, timezone) VALUES (?, ?, ?, ?, ?, ?)`
	result, execErr := a.ExecContext(ctx, tx, "insert", query, entity.Name, entity.Scope, entity.State, entity.City, entity.IsActive, entity.Timezone)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.holiday.create_calendar.exec_error", "name", entity.Name, "err", execErr)
		return 0, fmt.Errorf("insert holiday calendar: %w", execErr)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.create_calendar.last_id_error", "name", entity.Name, "err", err)
		return 0, fmt.Errorf("holiday calendar last insert id: %w", err)
	}

	calendar.SetID(uint64(id))
	return uint64(id), nil
}
