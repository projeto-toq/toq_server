package userentity

import (
	"database/sql"
)

// UserRoleEntity represents a row from the user_roles table in the database
//
// This entity associates users with roles and manages role-specific metadata.
// Supports role activation, status tracking, expiration, and temporary blocks.
//
// Schema Mapping:
//   - Database: user_roles table (InnoDB, utf8mb3)
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Foreign Keys: user_id → users.id, role_id → roles.id (CASCADE DELETE)
//   - Composite Unique Index: uk_user_roles on (user_id, role_id)
//   - Indexes: idx_user_roles_user, idx_user_roles_role, idx_user_roles_active, idx_user_roles_expires
//
// Table Purpose:
//   - Assign roles to users (one user can have multiple roles)
//   - Track active role for user session context
//   - Support role expiration (e.g., trial periods)
//   - Enable temporary role suspension (blocked_until)
//   - Manage role lifecycle and approval workflows
//
// NULL Handling:
//   - sql.NullTime: Used for expires_at and blocked_until (optional timestamps)
//   - Direct types: Used for NOT NULL columns
//
// Conversion:
//   - To Domain: Use permissionconverters.UserRoleEntityToDomain()
//   - From Domain: Use permissionconverters.UserRoleDomainToEntity()
//
// Business Rules (enforced by service layer):
//   - User can have multiple roles but only one active at a time
//   - Expired roles (expires_at < NOW) are automatically deactivated by cron job
//   - Blocked roles (blocked_until > NOW) cannot be activated
//   - Status codes: 0=pending, 1=approved, 2=rejected, 3=suspended
//   - CASCADE DELETE: user_roles removed when user or role is deleted
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type UserRoleEntity struct {
	// ID is the user-role association's unique identifier (PRIMARY KEY, AUTO_INCREMENT, INT UNSIGNED)
	ID uint32 `db:"id"`

	// UserID is the user's identifier (NOT NULL, INT UNSIGNED, FOREIGN KEY to users.id)
	// CASCADE DELETE: association removed when user deleted
	UserID uint32 `db:"user_id"`

	// RoleID is the role's identifier (NOT NULL, INT UNSIGNED, FOREIGN KEY to roles.id)
	// CASCADE DELETE: association removed when role deleted
	// References roles.id
	RoleID uint32 `db:"role_id"`

	// IsActive indicates if this is the user's currently active role (NOT NULL, TINYINT UNSIGNED, DEFAULT 1)
	// true = active (used for JWT claims and permission checks)
	// false = inactive (available but not currently selected)
	// Only ONE user_role per user should have is_active=1 at any time
	IsActive bool `db:"is_active"`

	// Status represents the approval/lifecycle state of the role assignment (NOT NULL, TINYINT signed, DEFAULT 0)
	// 0 = pending approval (awaiting admin review)
	// 1 = approved (active and operational)
	// 2 = rejected (denied by admin)
	// 3 = suspended (temporarily disabled)
	// Signed TINYINT allows for future negative states if needed
	Status int8 `db:"status"`

	// ExpiresAt is the optional expiration timestamp for the role (NULL, TIMESTAMP(6))
	// Used for trial periods, temporary access, or time-limited promotions
	// After this time, role should be deactivated by cron job
	// NULL = no expiration (permanent assignment)
	ExpiresAt sql.NullTime `db:"expires_at"`

	// BlockedUntil is the optional timestamp until which role is blocked (NULL, DATETIME)
	// Used for temporary suspensions (e.g., policy violation cooldown)
	// Role cannot be activated while blocked_until > NOW()
	// NULL = not blocked
	// Example: User banned for 7 days → blocked_until = NOW() + 7 days
	BlockedUntil sql.NullTime `db:"blocked_until"`
}
