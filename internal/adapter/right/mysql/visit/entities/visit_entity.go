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
//   - Foreign Keys: listing_id → listings(id), owner_id → users(id), realtor_id → users(id)
//   - Indexes: idx_visit_listing, idx_visit_scheduled_start (recommended)
//   - Status: ENUM('PENDING_OWNER', 'CONFIRMED', 'CANCELLED', 'DONE')
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

	// ListingID is the foreign key to the listings table (NOT NULL, INT UNSIGNED)
	// References the listing being visited
	ListingID int64

	// OwnerID is the foreign key to the users table (NOT NULL, INT UNSIGNED)
	// References the property owner who receives the visit request
	OwnerID int64

	// RealtorID is the foreign key to the users table (NOT NULL, INT UNSIGNED)
	// References the realtor who requests/conducts the visit
	RealtorID int64

	// ScheduledStart is the visit start date/time (NOT NULL, DATETIME)
	// Stored in UTC, converted to America/Sao_Paulo in service layer
	ScheduledStart time.Time

	// ScheduledEnd is the visit end date/time (NOT NULL, DATETIME)
	// Stored in UTC, must be after ScheduledStart (validated in service layer)
	ScheduledEnd time.Time

	// Status is the visit workflow state (NOT NULL, ENUM)
	// Allowed values: 'PENDING_OWNER', 'CONFIRMED', 'CANCELLED', 'DONE'
	// State transitions validated in service layer
	Status string

	// CancelReason explains why the visit was cancelled (NULL, VARCHAR(255))
	// Required when Status='CANCELLED', NULL otherwise
	CancelReason sql.NullString

	// Notes contains additional information about the visit (NULL, TEXT)
	// Optional field for extra context or instructions
	Notes sql.NullString

	// CreatedBy is the user ID who created this visit record (NOT NULL, INT UNSIGNED)
	// Audit field for tracking record origin
	CreatedBy int64

	// UpdatedBy is the user ID who last updated this visit (NULL, INT UNSIGNED)
	// Audit field for tracking record modifications
	UpdatedBy sql.NullInt64
}
