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

const calendarsMaxPageSize = 50

func (a *HolidayAdapter) ListCalendars(ctx context.Context, tx *sql.Tx, filter holidaymodel.CalendarListFilter) (holidaymodel.CalendarListResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return holidaymodel.CalendarListResult{}, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := make([]string, 0)
	args := make([]any, 0)

	if filter.Scope != nil {
		conditions = append(conditions, "scope = ?")
		args = append(args, string(*filter.Scope))
	}
	if filter.State != nil {
		conditions = append(conditions, "state = ?")
		args = append(args, *filter.State)
	}
	if filter.City != nil {
		conditions = append(conditions, "city = ?")
		args = append(args, *filter.City)
	}
	if filter.OnlyActive != nil {
		conditions = append(conditions, "is_active = ?")
		if *filter.OnlyActive {
			args = append(args, true)
		} else {
			args = append(args, false)
		}
	}
	if filter.SearchTerm != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+filter.SearchTerm+"%")
	}

	where := "1=1"
	if len(conditions) > 0 {
		where = strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM holiday_calendars WHERE %s", where)
	var total int64
	countRow := a.QueryRowContext(ctx, tx, "select", countQuery, args...)
	if err = countRow.Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.list_calendars.count_error", "err", err)
		return holidaymodel.CalendarListResult{}, fmt.Errorf("count holiday calendars: %w", err)
	}

	limit, offset := defaultPagination(filter.Limit, filter.Page, calendarsMaxPageSize)

	query := fmt.Sprintf(`
		SELECT id, name, scope, state, city, is_active, timezone
		FROM holiday_calendars
		WHERE %s
		ORDER BY name ASC
		LIMIT ? OFFSET ?
	`, where)

	listArgs := append(append([]any{}, args...), limit, offset)
	rows, queryErr := a.QueryContext(ctx, tx, "select", query, listArgs...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.holiday.list_calendars.query_error", "err", queryErr)
		return holidaymodel.CalendarListResult{}, fmt.Errorf("query holiday calendars: %w", queryErr)
	}
	defer rows.Close()

	calendars := make([]holidaymodel.CalendarInterface, 0)
	for rows.Next() {
		var calendarEntity entity.CalendarEntity
		if err = rows.Scan(&calendarEntity.ID, &calendarEntity.Name, &calendarEntity.Scope, &calendarEntity.State, &calendarEntity.City, &calendarEntity.IsActive, &calendarEntity.Timezone); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.holiday.list_calendars.scan_error", "err", err)
			return holidaymodel.CalendarListResult{}, fmt.Errorf("scan holiday calendar: %w", err)
		}
		calendars = append(calendars, converters.ToCalendarModel(calendarEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.list_calendars.rows_error", "err", err)
		return holidaymodel.CalendarListResult{}, fmt.Errorf("iterate holiday calendars: %w", err)
	}

	return holidaymodel.CalendarListResult{Calendars: calendars, Total: total}, nil
}
