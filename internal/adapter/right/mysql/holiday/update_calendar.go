package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/converters"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *HolidayAdapter) UpdateCalendar(ctx context.Context, tx *sql.Tx, calendar holidaymodel.CalendarInterface) error {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToCalendarEntity(calendar)

	query := `UPDATE holiday_calendars SET name = ?, scope = ?, state = ?, city_ibge = ?, is_active = ?, timezone = ? WHERE id = ?`
	result, err := exec.ExecContext(ctx, query, entity.Name, entity.Scope, entity.State, entity.CityIBGE, entity.IsActive, entity.Timezone, entity.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.update_calendar.exec_error", "calendar_id", entity.ID, "err", err)
		return fmt.Errorf("update holiday calendar: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.update_calendar.rows_error", "calendar_id", entity.ID, "err", err)
		return fmt.Errorf("holiday calendar rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
