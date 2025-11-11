package userentity

import (
	"database/sql"
	"time"
)

// UserWithRoleEntity represents the result of a JOIN query between users, user_roles and roles tables.
//
// This entity is used to efficiently load a user with their active role in a single database query,
// avoiding N+1 query problems. It combines data from three tables:
//   - users: user personal and account information
//   - user_roles: user-role association with status and expiration
//   - roles: role definition with slug, name and permissions
//
// Schema Mapping:
//   - Database: users LEFT JOIN user_roles ON users.id = user_roles.user_id
//     LEFT JOIN roles ON user_roles.role_id = roles.id
//   - Filter: user_roles.is_active = 1 (active role only)
//
// NULL Handling:
//   - All user_role and role fields are nullable (sql.Null*) to handle users without active roles
//   - Service layer decides if missing active role is an error or valid state
//
// Conversion:
//   - To Domain: Use userconverters.UserWithRoleEntityToDomain()
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - Prefer this over separate User + UserRole queries for performance
type UserWithRoleEntity struct {
	// ==================== User fields (from users table) ====================

	// UserID is the user's unique identifier (users.id, INT UNSIGNED)
	UserID uint32

	// FullName is the user's complete legal name (users.full_name, VARCHAR(150), NOT NULL)
	FullName string

	// NickName is the user's display name (users.nick_name, VARCHAR(45), NULL)
	NickName sql.NullString

	// NationalID is the user's CPF or CNPJ (users.national_id, VARCHAR(25), NOT NULL)
	NationalID string

	// CreciNumber is the CRECI registration number (users.creci_number, VARCHAR(15), NULL)
	CreciNumber sql.NullString

	// CreciState is the Brazilian state where CRECI is registered (users.creci_state, VARCHAR(2), NULL)
	CreciState sql.NullString

	// CreciValidity is the CRECI expiration date (users.creci_validity, DATE, NULL)
	CreciValidity sql.NullTime

	// BornAt is the user's date of birth (users.born_at, DATE, NOT NULL)
	BornAt time.Time

	// PhoneNumber is the user's mobile in E.164 format (users.phone_number, VARCHAR(25), NOT NULL)
	PhoneNumber string

	// Email is the user's email address (users.email, VARCHAR(45), NOT NULL)
	Email string

	// ZipCode is the Brazilian postal code (CEP) (users.zip_code, VARCHAR(8), NOT NULL)
	ZipCode string

	// Street is the street name (users.street, VARCHAR(150), NOT NULL)
	Street string

	// Number is the building/property number (users.number, VARCHAR(15), NOT NULL)
	Number string

	// Complement provides additional address info (users.complement, VARCHAR(150), NULL)
	Complement sql.NullString

	// Neighborhood is the district/neighborhood name (users.neighborhood, VARCHAR(150), NOT NULL)
	Neighborhood string

	// City is the city name (users.city, VARCHAR(150), NOT NULL)
	City string

	// State is the Brazilian state code (users.state, VARCHAR(2), NOT NULL)
	State string

	// Password is the bcrypt hash of user's password (users.password, VARCHAR(100), NOT NULL)
	Password string

	// OptStatus indicates if user opted in for marketing (users.opt_status, TINYINT UNSIGNED, NOT NULL)
	OptStatus bool

	// LastActivityAt is the timestamp of user's last action (users.last_activity_at, TIMESTAMP(6), NOT NULL)
	LastActivityAt time.Time

	// Deleted indicates soft delete status (users.deleted, TINYINT UNSIGNED, NOT NULL)
	Deleted bool

	// ==================== UserRole fields (from user_roles table) ====================
	// All nullable because user might not have active role (LEFT JOIN)

	// UserRoleID is the user-role association's unique identifier (user_roles.id, INT UNSIGNED, NULL)
	UserRoleID sql.NullInt32

	// UserRoleUserID is redundant user_id from JOIN (user_roles.user_id, INT UNSIGNED, NULL)
	UserRoleUserID sql.NullInt32

	// UserRoleRoleID is the role's identifier (user_roles.role_id, INT UNSIGNED, NULL)
	UserRoleRoleID sql.NullInt32

	// UserRoleIsActive indicates if this is the currently active role (user_roles.is_active, TINYINT UNSIGNED, NULL)
	UserRoleIsActive sql.NullBool

	// UserRoleStatus is the approval/lifecycle state (user_roles.status, TINYINT signed, NULL)
	// Uses sql.NullInt16 to support TINYINT signed (-128 to 127)
	UserRoleStatus sql.NullInt16

	// UserRoleExpiresAt is the optional expiration timestamp (user_roles.expires_at, TIMESTAMP(6), NULL)
	UserRoleExpiresAt sql.NullTime

	// UserRoleBlockedUntil is the optional block timestamp (user_roles.blocked_until, DATETIME, NULL)
	UserRoleBlockedUntil sql.NullTime

	// ==================== Role fields (from roles table) ====================
	// All nullable because user might not have active role (LEFT JOIN)

	// RoleID is the role's unique identifier (roles.id, INT UNSIGNED, NULL)
	RoleID sql.NullInt32

	// RoleSlug is the role's system identifier (roles.slug, VARCHAR(100), NULL)
	RoleSlug sql.NullString

	// RoleName is the role's display name (roles.name, VARCHAR(100), NULL)
	RoleName sql.NullString

	// RoleDescription is the role's detailed description (roles.description, TEXT, NULL)
	RoleDescription sql.NullString

	// RoleIsSystemRole indicates if role is system-defined (roles.is_system_role, TINYINT, NULL)
	RoleIsSystemRole sql.NullBool

	// RoleIsActive indicates if role is enabled (roles.is_active, TINYINT, NULL)
	RoleIsActive sql.NullBool
}
