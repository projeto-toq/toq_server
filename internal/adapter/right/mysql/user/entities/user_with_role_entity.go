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
	UserID            int64
	FullName          string
	NickName          sql.NullString
	NationalID        string
	CreciNumber       sql.NullString
	CreciState        sql.NullString
	CreciValidity     sql.NullTime
	BornAT            time.Time
	PhoneNumber       string
	Email             string
	ZipCode           string
	Street            string
	Number            string
	Complement        sql.NullString
	Neighborhood      string
	City              string
	State             string
	Password          string
	OptStatus         bool
	LastActivityAT    time.Time
	Deleted           bool
	LastSignInAttempt sql.NullTime

	// ==================== UserRole fields (from user_roles table) ====================
	// All nullable because user might not have active role
	UserRoleID           sql.NullInt64 // user_roles.id
	UserRoleUserID       sql.NullInt64 // user_roles.user_id (redundant but from JOIN)
	UserRoleRoleID       sql.NullInt64 // user_roles.role_id
	UserRoleIsActive     sql.NullBool  // user_roles.is_active
	UserRoleStatus       sql.NullInt64 // user_roles.status
	UserRoleExpiresAt    sql.NullTime  // user_roles.expires_at
	UserRoleBlockedUntil sql.NullTime  // user_roles.blocked_until

	// ==================== Role fields (from roles table) ====================
	// All nullable because user might not have active role
	RoleID           sql.NullInt64  // roles.id
	RoleSlug         sql.NullString // roles.slug
	RoleName         sql.NullString // roles.name
	RoleDescription  sql.NullString // roles.description
	RoleIsSystemRole sql.NullBool   // roles.is_system_role
	RoleIsActive     sql.NullBool   // roles.is_active
}
