package entities

import (
	"database/sql"
	"time"
)

// VisitEntity represents a row from the listing_visits table in the database.
//
// Schema mapping (InnoDB, utf8mb4_unicode_ci):
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Foreign Keys: listing_identity_id → listing_identities(id), user_id → users(id)
//   - Indexes: fk_visits_listing_identity_idx, fk_visits_user_idx, idx_scheduled_date, idx_status
//   - Status: ENUM('PENDING','APPROVED','REJECTED','CANCELLED','COMPLETED','NO_SHOW')
//   - Source: ENUM('APP','WEB','ADMIN') DEFAULT 'APP'
//   - requested_at: DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP (owner response metrics)
//
// NULL handling:
//   - sql.NullString for nullable TEXT/VARCHAR
//   - sql.NullTime for nullable DATETIME
//   - Direct types for NOT NULL columns
//
// Usage rules:
//   - Adapter layer only; convert via converters.ToVisitEntity/ToVisitModel
//   - Keep field order aligned with SELECT statements in visit adapter functions
type VisitEntity struct {
	// ID is the visit's unique identifier (PRIMARY KEY, AUTO_INCREMENT)
	ID int64

	// ListingIdentityID is the foreign key to the listing_identities table (NOT NULL, INT UNSIGNED)
	// References the listing identity being visited (all versions)
	ListingIdentityID int64

	// ListingVersion stores the listing version visible when the visit was requested (NOT NULL, TINYINT UNSIGNED)
	ListingVersion uint8

	// RequesterUserID is the foreign key to the users table (NOT NULL, INT UNSIGNED)
	// References the user (realtor/buyer) who requests the visit
	RequesterUserID int64

	// OwnerUserID is derived from listing_identities.user_id (not persisted in listing_visits)
	OwnerUserID int64

	// ScheduledStart is the visit start date/time (NOT NULL, DATETIME)
	// Stored in UTC, converted to America/Sao_Paulo in service layer
	ScheduledStart time.Time
	ScheduledEnd   time.Time

	// Status is the visit workflow state (NOT NULL, ENUM)
	// Allowed values: 'PENDING', 'APPROVED', 'REJECTED', 'CANCELLED', 'COMPLETED', 'NO_SHOW'
	// State transitions validated in service layer
	Status string

	// Source identifies where the visit was created (e.g., APP, WEB)
	Source sql.NullString

	// Notes stores requester/owner notes in a single field (TEXT)
	Notes sql.NullString

	// RejectionReason stores owner rejection reason or cancellation context
	RejectionReason sql.NullString

	FirstOwnerActionAt sql.NullTime

	// RequestedAt stores when the visit was requested; used for owner response time metrics (non-audit).
	// DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	RequestedAt time.Time
}
