package entity

import "database/sql"

// CalendarEntity mirrors the holiday_calendars table.
type CalendarEntity struct {
	ID       uint64
	Name     string
	Scope    string
	State    sql.NullString
	CityIBGE sql.NullString
	IsActive bool
	Timezone string
}
