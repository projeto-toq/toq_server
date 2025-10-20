package entity

import "time"

// DateEntity mirrors the holiday_calendar_dates table.
type DateEntity struct {
	ID         uint64
	CalendarID uint64
	Holiday    time.Time
	Label      string
	Recurrent  bool
}
