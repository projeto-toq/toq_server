package holidayservices

import (
	"context"
	"time"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	holidayrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/holiday_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

// HolidayServiceInterface orchestrates operations on holiday calendars.
type HolidayServiceInterface interface {
	CreateCalendar(ctx context.Context, input CreateCalendarInput) (holidaymodel.CalendarInterface, error)
	UpdateCalendar(ctx context.Context, input UpdateCalendarInput) (holidaymodel.CalendarInterface, error)
	GetCalendarByID(ctx context.Context, id uint64) (holidaymodel.CalendarInterface, error)
	ListCalendars(ctx context.Context, filter holidaymodel.CalendarListFilter) (holidaymodel.CalendarListResult, error)
	CreateCalendarDate(ctx context.Context, input CreateCalendarDateInput) (holidaymodel.CalendarDateInterface, error)
	UpdateCalendarDate(ctx context.Context, input UpdateCalendarDateInput) (holidaymodel.CalendarDateInterface, error)
	GetCalendarDateByID(ctx context.Context, id uint64, timezone string) (holidaymodel.CalendarDateInterface, error)
	DeleteCalendarDate(ctx context.Context, id uint64) error
	ListCalendarDates(ctx context.Context, filter holidaymodel.CalendarDatesFilter) (holidaymodel.CalendarDatesResult, error)
	// CleanOldCalendarDates deletes non-recurrent dates older than cutoff to enforce retention.
	CleanOldCalendarDates(ctx context.Context, cutoff time.Time, limit int) (int64, error)
}

type holidayService struct {
	repo          holidayrepository.HolidayRepositoryInterface
	globalService globalservice.GlobalServiceInterface
}

// NewHolidayService builds a new holiday service instance.
func NewHolidayService(
	repo holidayrepository.HolidayRepositoryInterface,
	globalService globalservice.GlobalServiceInterface,
) HolidayServiceInterface {
	return &holidayService{
		repo:          repo,
		globalService: globalService,
	}
}
