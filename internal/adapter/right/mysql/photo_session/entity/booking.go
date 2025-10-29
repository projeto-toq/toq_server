package entity

import (
	"database/sql"
	"time"
)

// Booking represents a row from photographer_photo_session_bookings.
type Booking struct {
	ID             uint64
	AgendaEntryID  uint64
	PhotographerID uint64
	ListingID      int64
	StartsAt       time.Time
	EndsAt         time.Time
	Status         string
	Reason         sql.NullString
}
