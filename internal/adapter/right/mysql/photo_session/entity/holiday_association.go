package entity

import "database/sql"

// HolidayAssociation represents a row from photographer_holiday_calendars.
type HolidayAssociation struct {
	ID                 uint64
	PhotographerUserID uint64
	HolidayCalendarID  uint64
	CreatedAt          sql.NullTime
}
