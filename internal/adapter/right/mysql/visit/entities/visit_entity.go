package entities

import (
	"database/sql"
	"time"
)

// VisitEntity represents a row from the listing_visits table in the database
//
// This struct maps directly to the database schema for visit scheduling and tracking.
// It uses sql.Null* types for nullable columns and should ONLY be used within
// the MySQL adapter layer.
//
// Schema Mapping:
//   - Database: listing_visits table (InnoDB, utf8mb4_unicode_ci)
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Foreign Keys: listing_identity_id → listing_identities(id), user_id → users(id)
//   - Indexes: idx_visit_listing, idx_visit_scheduled_start (recommended)
//   - Status: ENUM('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED', 'COMPLETED', 'NO_SHOW')
//
// NULL Handling:
//   - sql.NullString: Used for VARCHAR/TEXT columns that allow NULL
//   - sql.NullInt64: Used for INT columns that allow NULL (UpdatedBy)
//   - Direct types: Used for NOT NULL columns
//
// Conversion:
//   - To Domain: Use converters.ToVisitModel()
//   - From Domain: Use converters.ToVisitEntity()
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
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

	// OwnerUserID is the listing owner (NOT NULL, INT UNSIGNED)
	OwnerUserID int64

	// ScheduledStart is the visit start date/time (NOT NULL, DATETIME)
	// Stored in UTC, converted to America/Sao_Paulo in service layer
	ScheduledStart time.Time

	// ScheduledEnd is the visit end date/time (NOT NULL, DATETIME)
	// Stored in UTC, must be after ScheduledStart (validated in service layer)
	ScheduledEnd time.Time

	// DurationMinutes is the visit duration in minutes (NOT NULL, SMALLINT)
	DurationMinutes int64

	// Status is the visit workflow state (NOT NULL, ENUM)
	// Allowed values: 'PENDING', 'APPROVED', 'REJECTED', 'CANCELLED', 'COMPLETED', 'NO_SHOW'
	// State transitions validated in service layer
	Status string

	// Type indicates the kind of visit (NOT NULL, ENUM)
	// Allowed values: WITH_CLIENT, REALTOR_ONLY, CONTENT_PRODUCTION
	Type string

	// Source identifies where the visit was created (e.g., APP, WEB)
	Source sql.NullString

	// RealtorNotes stores notes from the requester/realtor
	RealtorNotes sql.NullString

	// OwnerNotes stores notes from the owner
	OwnerNotes sql.NullString

	// RejectionReason stores owner rejection reason
	RejectionReason sql.NullString

	// CancelReason explains why the visit was cancelled (NULL, VARCHAR(255))
	// Required when Status='CANCELLED', NULL otherwise
	CancelReason sql.NullString

	// FirstOwnerActionAt stores when the owner first approved/rejected the visit (NULL, DATETIME)
	FirstOwnerActionAt sql.NullTime
}
