package holidaymodel

import "time"

// CalendarDateInterface describes a single holiday date.
type CalendarDateInterface interface {
	ID() uint64
	SetID(id uint64)
	CalendarID() uint64
	SetCalendarID(id uint64)
	HolidayDate() time.Time
	SetHolidayDate(value time.Time)
	Label() string
	SetLabel(value string)
	IsRecurrent() bool
	SetRecurrent(value bool)
}

// NewCalendarDate builds an empty calendar date entity.
func NewCalendarDate() CalendarDateInterface {
	return &calendarDate{}
}
