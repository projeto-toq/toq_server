package entity

import (
	"database/sql"
	"time"
)

// EntryEntity mirrors the listing_agenda_entries table.
type EntryEntity struct {
	ID             uint64
	AgendaID       uint64
	EntryType      string
	StartsAt       time.Time
	EndsAt         time.Time
	Blocking       bool
	Reason         sql.NullString
	VisitID        sql.NullInt64
	PhotoBookingID sql.NullInt64
	CreatedBy      int64
	UpdatedBy      sql.NullInt64
}
