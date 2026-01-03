package scheduleentity

import (
	"database/sql"
	"time"
)

// EntryEntity maps listing_agenda_entries rows (time-bounded slots) for adapter use only.
type EntryEntity struct {
	ID             uint64         // PK AUTO_INCREMENT
	AgendaID       uint64         // FK listing_agendas.id (INT UNSIGNED NOT NULL)
	EntryType      string         // ENUM('BLOCK','TEMP_BLOCK','VISIT_PENDING','VISIT_CONFIRMED','PHOTO_SESSION','HOLIDAY_INFO')
	StartsAt       time.Time      // DATETIME NOT NULL
	EndsAt         time.Time      // DATETIME NOT NULL
	Blocking       bool           // TINYINT(1) NOT NULL
	Reason         sql.NullString // VARCHAR(120) NULL
	VisitID        sql.NullInt64  // FK listing_visits.id NULL
	PhotoBookingID sql.NullInt64  // FK photographer_photo_session_bookings.id NULL
}
