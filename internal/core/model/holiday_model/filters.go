package holidaymodel

import "time"

// CalendarListFilter constrains calendar lookups for admin endpoints.
type CalendarListFilter struct {
	Scope      *CalendarScope
	State      *string
	CityIBGE   *string
	SearchTerm string
	OnlyActive *bool
	Page       int
	Limit      int
}

// CalendarListResult represents a paginated set of calendars.
type CalendarListResult struct {
	Calendars []CalendarInterface
	Total     int64
}

// CalendarDatesFilter constrains date lookups for a calendar.
type CalendarDatesFilter struct {
	CalendarID uint64
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

// CalendarDatesResult represents holiday dates paginated.
type CalendarDatesResult struct {
	Dates []CalendarDateInterface
	Total int64
}
