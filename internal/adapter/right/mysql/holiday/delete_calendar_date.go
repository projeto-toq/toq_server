package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *HolidayAdapter) DeleteCalendarDate(ctx context.Context, tx *sql.Tx, id uint64) error {
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

	query := `DELETE FROM holiday_calendar_dates WHERE id = ?`
	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.delete_date.exec_error", "date_id", id, "err", err)
		return fmt.Errorf("delete holiday calendar date: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.delete_date.rows_error", "date_id", id, "err", err)
		return fmt.Errorf("holiday calendar date rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
