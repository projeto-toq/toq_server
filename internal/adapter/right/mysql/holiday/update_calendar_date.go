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

	entity := converters.ToDateEntity(date)

	query := `UPDATE holiday_calendar_dates SET holiday_date = ?, label = ?, is_recurrent = ? WHERE id = ?`
	defer a.ObserveOnComplete("update", query)()
	result, err := exec.ExecContext(ctx, query, entity.Holiday, entity.Label, entity.Recurrent, entity.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.update_calendar_date.exec_error", "date_id", entity.ID, "err", err)
		return fmt.Errorf("update holiday calendar date: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.update_calendar_date.rows_error", "date_id", entity.ID, "err", err)
		return fmt.Errorf("holiday calendar date rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
