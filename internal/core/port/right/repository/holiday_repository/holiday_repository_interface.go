package holidayrepository

import (
	"context"
	"database/sql"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
)

// HolidayRepositoryInterface defines persistence operations for holiday calendars.
type HolidayRepositoryInterface interface {
	CreateCalendar(ctx context.Context, tx *sql.Tx, calendar holidaymodel.CalendarInterface) (uint64, error)
	UpdateCalendar(ctx context.Context, tx *sql.Tx, calendar holidaymodel.CalendarInterface) error
	GetCalendarByID(ctx context.Context, tx *sql.Tx, id uint64) (holidaymodel.CalendarInterface, error)
	ListCalendars(ctx context.Context, tx *sql.Tx, filter holidaymodel.CalendarListFilter) (holidaymodel.CalendarListResult, error)
	CreateCalendarDate(ctx context.Context, tx *sql.Tx, date holidaymodel.CalendarDateInterface) (uint64, error)
	DeleteCalendarDate(ctx context.Context, tx *sql.Tx, id uint64) error
	ListCalendarDates(ctx context.Context, tx *sql.Tx, filter holidaymodel.CalendarDatesFilter) (holidaymodel.CalendarDatesResult, error)
}
