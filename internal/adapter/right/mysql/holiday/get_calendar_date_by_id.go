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

// GetCalendarDateByID fetches a holiday calendar date by its identifier.
func (a *HolidayAdapter) GetCalendarDateByID(ctx context.Context, tx *sql.Tx, id uint64) (holidaymodel.CalendarDateInterface, error) {
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

	query := `SELECT id, calendar_id, holiday_date, label, is_recurrent FROM holiday_calendar_dates WHERE id = ?`
	row := exec.QueryRowContext(ctx, query, id)

	var dateEntity entity.DateEntity
	if err = row.Scan(&dateEntity.ID, &dateEntity.CalendarID, &dateEntity.Holiday, &dateEntity.Label, &dateEntity.Recurrent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.holiday.get_calendar_date.scan_error", "date_id", id, "err", err)
		return nil, fmt.Errorf("scan holiday calendar date: %w", err)
	}

	return converters.ToDateModel(dateEntity), nil
}
