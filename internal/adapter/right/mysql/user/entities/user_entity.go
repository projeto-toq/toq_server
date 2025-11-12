package userentity

import (
	"database/sql"
	"time"
)

// UserEntity represents a row from the users table in the database
//
// This struct maps directly to the database schema and uses sql.Null* types
// for nullable columns. It should ONLY be used within the MySQL adapter layer.
//
// Schema Mapping:
//   - Database: users table (InnoDB, utf8mb3)
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Unique Constraints: national_id, email, phone_number
//   - Indexes: idx_users_deleted, idx_users_national_id, idx_users_email
//
// NULL Handling:
//   - sql.NullString: Used for VARCHAR columns that allow NULL
//   - sql.NullTime: Used for TIMESTAMP(6)/DATE columns that allow NULL
//   - Direct types: Used for NOT NULL columns
//
// Conversion:
//   - To Domain: Use userconverters.UserEntityToDomain()
//   - From Domain: Use userconverters.UserDomainToEntity()
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type UserEntity struct {
	// ID is the user's unique identifier (PRIMARY KEY, AUTO_INCREMENT, INT UNSIGNED)
	ID uint32

	// FullName is the user's complete legal name (NOT NULL, VARCHAR(150))
	// Used for legal identification and contracts
	// Example: "João Silva Santos"
	FullName string

	// NickName is the user's display name (NULL, VARCHAR(45))
	// Used in UI elements and notifications
	// Example: "João"
	NickName sql.NullString

	// NationalID is the user's CPF or CNPJ (NOT NULL, VARCHAR(25), UNIQUE)
	// Format: digits only (no punctuation)
	// CPF example: "12345678901" (11 digits)
	// CNPJ example: "12345678000195" (14 digits)
	NationalID string

	// CreciNumber is the CRECI registration number (NULL, VARCHAR(15))
	// Required ONLY for realtor role
	// Format: numeric followed by "-F" (e.g., "12345-F")
	CreciNumber sql.NullString

	// CreciState is the Brazilian state where CRECI is registered (NULL, VARCHAR(2))
	// Required when CreciNumber is provided
	// Must be a valid 2-letter Brazilian state code (e.g., "SP", "RJ")
	CreciState sql.NullString

	// CreciValidity is the CRECI expiration date (NULL, DATE)
	// Must be a future date for active realtors
	CreciValidity sql.NullTime

	// BornAt is the user's date of birth (NOT NULL, DATE)
	// User must be at least 18 years old (enforced by service layer)
	BornAt time.Time

	// PhoneNumber is the user's mobile in E.164 format (NOT NULL, VARCHAR(25), UNIQUE)
	// Must include country code with + prefix
	// Example: "+5511999999999" (Brazil mobile)
	PhoneNumber string

	// Email is the user's email address (NOT NULL, VARCHAR(45), UNIQUE)
	// Used for account recovery and email notifications
	// Example: "joao.silva@example.com"
	Email string

	// ZipCode is the Brazilian postal code (CEP) (NOT NULL, VARCHAR(8))
	// Format: 8 digits without separators (no hyphen)
	// Example: "01310100" (Avenida Paulista, São Paulo)
	ZipCode string

	// Street is the street name (NOT NULL, VARCHAR(150))
	// Typically populated automatically from CEP lookup
	// Example: "Avenida Paulista"
	Street string

	// Number is the building/property number (NOT NULL, VARCHAR(15))
	// Use "S/N" for addresses without number
	Number string

	// Complement provides additional address info (NULL, VARCHAR(150))
	// Optional field for apartment number, building name, etc.
	// Example: "Apto 501", "Bloco B"
	Complement sql.NullString

	// Neighborhood is the district/neighborhood name (NOT NULL, VARCHAR(150))
	// Example: "Bela Vista"
	Neighborhood string

	// City is the city name (NOT NULL, VARCHAR(150))
	// Example: "São Paulo"
	City string

	// State is the Brazilian state code (NOT NULL, VARCHAR(2))
	// Must be a valid 2-letter state code
	// Example: "SP"
	State string

	// Password is the bcrypt hash of user's password (NOT NULL, VARCHAR(100))
	// NEVER store plain text passwords
	// Hash generated with bcrypt.GenerateFromPassword()
	Password string

	// OptStatus indicates if user opted in for marketing (NOT NULL, TINYINT UNSIGNED, DEFAULT 0)
	// true = opted in, false = opted out
	OptStatus bool

	// LastActivityAt is the timestamp of user's last action (NOT NULL, TIMESTAMP(6))
	// Updated by activity tracking system
	// Used for idle timeout and analytics
	LastActivityAt time.Time

	// Deleted indicates soft delete status (NOT NULL, TINYINT UNSIGNED, DEFAULT 0)
	// true = logically deleted (hidden from queries)
	// false = active user
	Deleted bool

	// ==================== NEW: User-level blocking fields ====================

	// BlockedUntil is the timestamp until which the user is temporarily blocked (NULL, DATETIME)
	// User is blocked while blocked_until > NOW() (temporal block)
	// NULL = not temporarily blocked
	// Used for security measures (failed signin attempts, suspicious activity)
	// Worker automatically clears when timestamp expires
	// Example: User banned for 15 minutes → blocked_until = NOW() + 15 minutes
	BlockedUntil sql.NullTime

	// PermanentlyBlocked indicates if user is permanently blocked by admin (NOT NULL, TINYINT(1))
	// 1 = permanently blocked (cannot authenticate, requires manual unblock)
	// 0 = not permanently blocked (normal status)
	// Used for policy violations, fraud, terms of service violations
	// Only admin can set/unset this flag
	// Example: User found committing fraud → permanently_blocked = 1
	PermanentlyBlocked bool
}
