package scheduleentity

import (
	"database/sql"
	"time"
)

// EntryEntity maps listing_agenda_entries rows (InnoDB, utf8mb4) and is restricted to the MySQL adapter layer.
// Schema summary: PK (id), FK agenda_id -> listing_agendas.id, optional visit_id/photo_booking_id, time-bounded slots.
// Nullable columns use sql.Null* to preserve DB semantics; conversions are handled by scheduleconverters.
type EntryEntity struct {
	// ID is the AUTO_INCREMENT primary key (INT UNSIGNED NOT NULL).
	ID uint64
	// AgendaID references listing_agendas.id (INT UNSIGNED NOT NULL, indexed fk_entries_agenda_idx).
	AgendaID uint64
	// EntryType stores the source of the entry (ENUM('BLOCK','TEMP_BLOCK','VISIT_PENDING','VISIT_CONFIRMED','PHOTO_SESSION','HOLIDAY_INFO') NOT NULL).
	EntryType string
	// StartsAt is the interval start (DATETIME NOT NULL, uses tz per agenda context in services).
	StartsAt time.Time
	// EndsAt is the interval end (DATETIME NOT NULL, must be greater than StartsAt enforced by services).
	EndsAt time.Time
	// Blocking flags if the slot blocks availability (TINYINT(1) NOT NULL).
	Blocking bool
	// Reason stores optional human-readable reason (VARCHAR(120) NULL).
	Reason sql.NullString
	// VisitID references listing_visits.id when the slot is tied to a visit (INT UNSIGNED NULL).
	VisitID sql.NullInt64
	// PhotoBookingID references photographer_photo_session_bookings.id when applicable (INT UNSIGNED NULL).
	PhotoBookingID sql.NullInt64
}
