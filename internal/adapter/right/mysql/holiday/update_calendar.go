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
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToCalendarEntity(calendar)

	query := `UPDATE holiday_calendars SET name = ?, scope = ?, state = ?, city = ?, is_active = ?, timezone = ? WHERE id = ?`
	result, execErr := a.ExecContext(ctx, tx, "update", query, entity.Name, entity.Scope, entity.State, entity.City, entity.IsActive, entity.Timezone, entity.ID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.holiday.update_calendar.exec_error", "calendar_id", entity.ID, "err", execErr)
		return fmt.Errorf("update holiday calendar: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.holiday.update_calendar.rows_error", "calendar_id", entity.ID, "err", rowsErr)
		return fmt.Errorf("holiday calendar rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
