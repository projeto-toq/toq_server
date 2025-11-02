package mysqlholidayadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/converters"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *HolidayAdapter) CreateCalendarDate(ctx context.Context, tx *sql.Tx, date holidaymodel.CalendarDateInterface) (uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToDateEntity(date)

	query := `INSERT INTO holiday_calendar_dates (calendar_id, holiday_date, label, is_recurrent) VALUES (?, ?, ?, ?)`
	result, execErr := a.ExecContext(ctx, tx, "insert", query, entity.CalendarID, entity.Holiday, entity.Label, entity.Recurrent)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.holiday.create_date.exec_error", "calendar_id", entity.CalendarID, "err", execErr)
		return 0, fmt.Errorf("insert holiday calendar date: %w", execErr)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.create_date.last_id_error", "calendar_id", entity.CalendarID, "err", err)
		return 0, fmt.Errorf("holiday calendar date last insert id: %w", err)
	}

	date.SetID(uint64(id))
	return uint64(id), nil
}
