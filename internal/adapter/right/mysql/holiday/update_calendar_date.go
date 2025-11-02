package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/converters"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateCalendarDate persists changes to an existing holiday calendar date.
func (a *HolidayAdapter) UpdateCalendarDate(ctx context.Context, tx *sql.Tx, date holidaymodel.CalendarDateInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToDateEntity(date)

	query := `UPDATE holiday_calendar_dates SET holiday_date = ?, label = ?, is_recurrent = ? WHERE id = ?`
	result, execErr := a.ExecContext(ctx, tx, "update", query, entity.Holiday, entity.Label, entity.Recurrent, entity.ID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.holiday.update_calendar_date.exec_error", "date_id", entity.ID, "err", execErr)
		return fmt.Errorf("update holiday calendar date: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.holiday.update_calendar_date.rows_error", "date_id", entity.ID, "err", rowsErr)
		return fmt.Errorf("holiday calendar date rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
