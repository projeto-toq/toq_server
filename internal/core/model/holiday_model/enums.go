package holidaymodel

// CalendarScope identifies the geographic scope a calendar applies to.
type CalendarScope string

const (
	ScopeNational CalendarScope = "NATIONAL"
	ScopeState    CalendarScope = "STATE"
	ScopeCity     CalendarScope = "CITY"
)
