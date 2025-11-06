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
//   - Database: users table (InnoDB, utf8mb4_unicode_ci)
//   - Primary Key: id (BIGINT AUTO_INCREMENT)
//   - Unique Constraints: national_id, email, phone_number
//   - Indexes: idx_users_deleted, idx_users_national_id, idx_users_email
//
// NULL Handling:
//   - sql.NullString: Used for VARCHAR columns that allow NULL
//   - sql.NullTime: Used for DATETIME columns that allow NULL
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
	// ID is the user's unique identifier (PRIMARY KEY, AUTO_INCREMENT, BIGINT)
	ID int64

	// FullName is the user's complete legal name (NOT NULL, VARCHAR(100))
	// Used for legal identification and contracts
	// Example: "João Silva Santos"
	FullName string

	// NickName is the user's display name (NULL, VARCHAR(50))
	// Used in UI elements and notifications
	// Example: "João"
	NickName sql.NullString

	// NationalID is the user's CPF or CNPJ (NOT NULL, VARCHAR(14), UNIQUE)
	// Format: digits only (no punctuation)
	// CPF example: "12345678901" (11 digits)
	// CNPJ example: "12345678000195" (14 digits)
	NationalID string

	// CreciNumber is the CRECI registration number (NULL, VARCHAR(20))
	// Required ONLY for realtor role
	// Format: numeric followed by "-F" (e.g., "12345-F")
	CreciNumber sql.NullString

	// CreciState is the Brazilian state where CRECI is registered (NULL, CHAR(2))
	// Required when CreciNumber is provided
	// Must be a valid 2-letter Brazilian state code (e.g., "SP", "RJ")
	CreciState sql.NullString

	// CreciValidity is the CRECI expiration date (NULL, DATE)
	// Must be a future date for active realtors
	CreciValidity sql.NullTime

	// BornAT is the user's date of birth (NOT NULL, DATE)
	// User must be at least 18 years old (enforced by service layer)
	BornAT time.Time

	// PhoneNumber is the user's mobile in E.164 format (NOT NULL, VARCHAR(20), UNIQUE)
	// Must include country code with + prefix
	// Example: "+5511999999999" (Brazil mobile)
	PhoneNumber string

	// Email is the user's email address (NOT NULL, VARCHAR(100), UNIQUE)
	// Used for account recovery and email notifications
	// Example: "joao.silva@example.com"
	Email string

	// ZipCode is the Brazilian postal code (CEP) (NOT NULL, CHAR(8))
	// Format: 8 digits without separators (no hyphen)
	// Example: "01310100" (Avenida Paulista, São Paulo)
	ZipCode string

	// Street is the street name (NOT NULL, VARCHAR(100))
	// Typically populated automatically from CEP lookup
	// Example: "Avenida Paulista"
	Street string

	// Number is the building/property number (NOT NULL, VARCHAR(10))
	// Use "S/N" for addresses without number
	Number string

	// Complement provides additional address info (NULL, VARCHAR(50))
	// Optional field for apartment number, building name, etc.
	// Example: "Apto 501", "Bloco B"
	Complement sql.NullString

	// Neighborhood is the district/neighborhood name (NOT NULL, VARCHAR(50))
	// Example: "Bela Vista"
	Neighborhood string

	// City is the city name (NOT NULL, VARCHAR(50))
	// Example: "São Paulo"
	City string

	// State is the Brazilian state code (NOT NULL, CHAR(2))
	// Must be a valid 2-letter state code
	// Example: "SP"
	State string

	// Password is the bcrypt hash of user's password (NOT NULL, VARCHAR(255))
	// NEVER store plain text passwords
	// Hash generated with bcrypt.GenerateFromPassword()
	Password string

	// OptStatus indicates if user opted in for marketing (NOT NULL, TINYINT(1), DEFAULT 0)
	// true = opted in, false = opted out
	OptStatus bool

	// LastActivityAT is the timestamp of user's last action (NOT NULL, DATETIME)
	// Updated by activity tracking system
	// Used for idle timeout and analytics
	LastActivityAT time.Time

	// Deleted indicates soft delete status (NOT NULL, TINYINT(1), DEFAULT 0)
	// true = logically deleted (hidden from queries)
	// false = active user
	Deleted bool

	// LastSignInAttempt is the timestamp of last sign-in attempt (NULL, DATETIME)
	// Used for security monitoring and wrong password tracking
	LastSignInAttempt sql.NullTime
}
