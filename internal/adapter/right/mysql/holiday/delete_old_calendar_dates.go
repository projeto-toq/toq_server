package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	holidayrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/holiday_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

var _ holidayrepository.HolidayRepositoryInterface = (*HolidayAdapter)(nil)

// DeleteOldCalendarDates removes non-recurrent holiday dates older than the cutoff.
// Returns rows deleted; zero rows is success (no data to prune).
func (a *HolidayAdapter) DeleteOldCalendarDates(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 {
		limit = 500
	}

	query := `DELETE FROM holiday_calendar_dates 
        WHERE is_recurrent = 0 AND holiday_date < ?
        LIMIT ?`

	res, execErr := a.ExecContext(ctx, tx, "delete_old_calendar_dates", query, cutoff, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.holiday.delete_old_calendar_dates.exec_error", "cutoff", cutoff, "limit", limit, "error", execErr)
		return 0, fmt.Errorf("delete old holiday dates: %w", execErr)
	}

	rows, raErr := res.RowsAffected()
	if raErr != nil {
		logger.Warn("mysql.holiday.delete_old_calendar_dates.rows_affected_warning", "error", raErr)
		return 0, nil
	}

	if rows > 0 {
		logger.Debug("mysql.holiday.delete_old_calendar_dates.success", "deleted", rows, "cutoff", cutoff, "limit", limit)
	}

	return rows, nil
}
