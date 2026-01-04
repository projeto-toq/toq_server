package entity

import "database/sql"

// Booking models photographer_photo_session_bookings with nullable fields preserved.
// Columns: id (PK, NOT NULL), photographer_user_id (NOT NULL), listing_identity_id (NOT NULL),
// agenda_entry_id (NOT NULL), starts_at (DATETIME(6) NOT NULL), ends_at (DATETIME(6) NOT NULL),
// status (ENUM NOT NULL), reason (VARCHAR(255) NULL), reservation_token (VARCHAR(36) NULL),
// reserved_until (DATETIME(6) NOT NULL DEFAULT DATE_ADD(CURRENT_TIMESTAMP(6), INTERVAL 3 DAY)).
type Booking struct {
	ID                uint64         // photographer_photo_session_bookings.id
	AgendaEntryID     uint64         // photographer_photo_session_bookings.agenda_entry_id (NOT NULL)
	PhotographerID    uint64         // photographer_photo_session_bookings.photographer_user_id (NOT NULL)
	ListingIdentityID int64          // photographer_photo_session_bookings.listing_identity_id (NOT NULL)
	StartsAt          sql.NullTime   // photographer_photo_session_bookings.starts_at (DATETIME(6), NOT NULL)
	EndsAt            sql.NullTime   // photographer_photo_session_bookings.ends_at (DATETIME(6), NOT NULL)
	Status            string         // photographer_photo_session_bookings.status (ENUM, NOT NULL)
	Reason            sql.NullString // photographer_photo_session_bookings.reason (NULLABLE)
	ReservationToken  sql.NullString // photographer_photo_session_bookings.reservation_token (NULLABLE)
	ReservedUntil     sql.NullTime   // photographer_photo_session_bookings.reserved_until (DATETIME(6), NOT NULL DEFAULT)
}
