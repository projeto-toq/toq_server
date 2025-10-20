package entity

import (
	"database/sql"
	"time"
)

// VisitEntity mirrors the listing_visits table.
type VisitEntity struct {
	ID             int64
	ListingID      int64
	OwnerID        int64
	RealtorID      int64
	ScheduledStart time.Time
	ScheduledEnd   time.Time
	Status         string
	CancelReason   sql.NullString
	Notes          sql.NullString
	CreatedBy      int64
	UpdatedBy      sql.NullInt64
}
