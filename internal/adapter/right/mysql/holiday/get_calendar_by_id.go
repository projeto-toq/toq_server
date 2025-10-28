package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/entity"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *HolidayAdapter) GetCalendarByID(ctx context.Context, tx *sql.Tx, id uint64) (holidaymodel.CalendarInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, name, scope, state, city, is_active, timezone FROM holiday_calendars WHERE id = ?`
	row := exec.QueryRowContext(ctx, query, id)

	var calendarEntity entity.CalendarEntity
	if err = row.Scan(&calendarEntity.ID, &calendarEntity.Name, &calendarEntity.Scope, &calendarEntity.State, &calendarEntity.City, &calendarEntity.IsActive, &calendarEntity.Timezone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.get_calendar.scan_error", "calendar_id", id, "err", err)
		return nil, fmt.Errorf("scan holiday calendar: %w", err)
	}

	return converters.ToCalendarModel(calendarEntity), nil
}
