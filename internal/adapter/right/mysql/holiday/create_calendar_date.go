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
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return 0, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToDateEntity(date)

	query := `INSERT INTO holiday_calendar_dates (calendar_id, holiday_date, label, is_recurrent) VALUES (?, ?, ?, ?)`
	defer a.ObserveOnComplete("insert", query)()
	result, err := exec.ExecContext(ctx, query, entity.CalendarID, entity.Holiday, entity.Label, entity.Recurrent)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.create_date.exec_error", "calendar_id", entity.CalendarID, "err", err)
		return 0, fmt.Errorf("insert holiday calendar date: %w", err)
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
