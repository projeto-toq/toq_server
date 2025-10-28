package holidayservices

import (
	"time"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
)

// CreateCalendarInput carries the data required to create a holiday calendar.
type CreateCalendarInput struct {
	Name     string
	Scope    holidaymodel.CalendarScope
	State    string
	City     string
	IsActive bool
	Timezone string
}

// UpdateCalendarInput captures the information to update an existing calendar.
type UpdateCalendarInput struct {
	ID       uint64
	Name     string
	Scope    holidaymodel.CalendarScope
	State    string
	City     string
	IsActive bool
	Timezone string
}

// CreateCalendarDateInput describes the payload to register a holiday date.
type CreateCalendarDateInput struct {
	CalendarID  uint64
	HolidayDate time.Time
	Label       string
	Recurrent   bool
}

// UpdateCalendarDateInput captures the data to update an existing holiday date.
type UpdateCalendarDateInput struct {
	ID          uint64
	CalendarID  uint64
	HolidayDate time.Time
	Label       string
	Recurrent   bool
}
