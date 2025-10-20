package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/entity"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const calendarDatesMaxPageSize = 100

func (a *HolidayAdapter) ListCalendarDates(ctx context.Context, tx *sql.Tx, filter holidaymodel.CalendarDatesFilter) (holidaymodel.CalendarDatesResult, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return holidaymodel.CalendarDatesResult{}, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := []string{"calendar_id = ?"}
	args := []any{filter.CalendarID}

	if filter.From != nil {
		conditions = append(conditions, "holiday_date >= ?")
		args = append(args, *filter.From)
	}
	if filter.To != nil {
		conditions = append(conditions, "holiday_date <= ?")
		args = append(args, *filter.To)
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM holiday_calendar_dates WHERE %s", where)
	var total int64
	if err = exec.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.list_dates.count_error", "calendar_id", filter.CalendarID, "err", err)
		return holidaymodel.CalendarDatesResult{}, fmt.Errorf("count holiday calendar dates: %w", err)
	}

	limit, offset := defaultPagination(filter.Limit, filter.Page, calendarDatesMaxPageSize)

	query := fmt.Sprintf(`
		SELECT id, calendar_id, holiday_date, label, is_recurrent
		FROM holiday_calendar_dates
		WHERE %s
		ORDER BY holiday_date ASC
		LIMIT ? OFFSET ?
	`, where)

	rows, err := exec.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.list_dates.query_error", "calendar_id", filter.CalendarID, "err", err)
		return holidaymodel.CalendarDatesResult{}, fmt.Errorf("query holiday calendar dates: %w", err)
	}
	defer rows.Close()

	dates := make([]holidaymodel.CalendarDateInterface, 0)
	for rows.Next() {
		var dateEntity entity.DateEntity
		if err = rows.Scan(&dateEntity.ID, &dateEntity.CalendarID, &dateEntity.Holiday, &dateEntity.Label, &dateEntity.Recurrent); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.holiday.list_dates.scan_error", "calendar_id", filter.CalendarID, "err", err)
			return holidaymodel.CalendarDatesResult{}, fmt.Errorf("scan holiday calendar date: %w", err)
		}
		dates = append(dates, converters.ToDateModel(dateEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.list_dates.rows_error", "calendar_id", filter.CalendarID, "err", err)
		return holidaymodel.CalendarDatesResult{}, fmt.Errorf("iterate holiday calendar dates: %w", err)
	}

	return holidaymodel.CalendarDatesResult{Dates: dates, Total: total}, nil
}
