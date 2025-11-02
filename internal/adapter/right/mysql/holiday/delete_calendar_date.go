package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *HolidayAdapter) DeleteCalendarDate(ctx context.Context, tx *sql.Tx, id uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM holiday_calendar_dates WHERE id = ?`
	result, execErr := a.ExecContext(ctx, tx, "delete", query, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.holiday.delete_date.exec_error", "date_id", id, "err", execErr)
		return fmt.Errorf("delete holiday calendar date: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.holiday.delete_date.rows_error", "date_id", id, "err", rowsErr)
		return fmt.Errorf("holiday calendar date rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
