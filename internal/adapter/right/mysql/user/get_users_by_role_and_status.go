package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUsersByRoleAndStatus returns users that have an active user_role with the given role slug and status
//
// This function retrieves users filtered by their current role assignment and status.
// Used for administrative queries and bulk operations on specific user segments.
//
// Query Logic:
//   - JOINs users with user_roles and roles tables
//   - Filters by ur.is_active = 1 (only current/active role assignment)
//   - Filters by ur.status (e.g., 0=pending, 1=active, 2=blocked)
//   - Filters by r.slug (role identifier like "owner", "realtor", "agency")
//   - Returns only non-deleted users (u.deleted = 0)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - role: Role slug to filter by (e.g., "owner", "realtor")
//   - status: UserRoleStatus to filter by (e.g., Active, Blocked, Pending)
//
// Returns:
//   - users: Slice of UserInterface matching criteria
//   - error: sql.ErrNoRows if no users found, or database errors
//
// Business Rules:
//   - Only returns users with ACTIVE role assignment (ur.is_active = 1)
//   - Users with multiple roles: only considered if specified role is active
//   - Soft-deleted users (deleted=1) are excluded
//   - Returns sql.ErrNoRows if no matches found (service maps to 404 or empty list)
//
// Edge Cases:
//   - User has multiple roles but queried role is inactive: NOT returned
//   - User has no roles: NOT returned (INNER JOIN filters them out)
//   - Invalid role slug: Returns sql.ErrNoRows (no matches)
//
// Performance:
//   - Uses indexes: user_roles.user_id, user_roles.is_active, roles.slug
//   - INNER JOIN ensures only users with matching role are processed
//
// Use Cases:
//   - Listing all active owners for bulk email
//   - Finding blocked realtors for administrative review
//   - Counting pending agencies for approval workflow
//
// Example:
//
//	users, err := adapter.GetUsersByRoleAndStatus(ctx, nil,
//	    permissionmodel.RoleSlugOwner,
//	    permissionmodel.UserRoleStatusActive)
//	if err == sql.ErrNoRows {
//	    // No active owners found
//	}
//	// users contains all active owners with deleted=0
func (ua *UserAdapter) GetUsersByRoleAndStatus(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleSlug, status globalmodel.UserRoleStatus) ([]usermodel.UserInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query joins users with roles and user_roles to filter by active role and status
	// Note: INNER JOIN excludes users without matching role
	// Note: ur.is_active = 1 ensures only current role assignment
	query := `
        SELECT u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state, u.creci_validity,
               u.born_at, u.phone_number, u.email, u.zip_code, u.street, u.number, u.complement,
               u.neighborhood, u.city, u.state, u.password, u.opt_status, u.last_activity_at, u.deleted, u.last_signin_attempt
          FROM users u
          JOIN user_roles ur ON ur.user_id = u.id AND ur.is_active = 1 AND ur.status = ?
          JOIN roles r ON r.id = ur.role_id AND r.slug = ?
         WHERE u.deleted = 0`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, int(status), role)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_users_by_role.query_error", "error", queryErr, "role", role, "status", status)
		return nil, fmt.Errorf("get users by role and status query: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entities
	entities, qerr := scanUserEntities(rows)
	if qerr != nil {
		utils.SetSpanError(ctx, qerr)
		logger.Error("mysql.user.get_users_by_role.scan_error", "error", qerr)
		return nil, fmt.Errorf("scan users by role rows: %w", qerr)
	}

	// Return sql.ErrNoRows if no users found (service layer handles mapping)
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Convert entities to domain models
	users := make([]usermodel.UserInterface, 0, len(entities))
	for _, e := range entities {
		u := userconverters.UserEntityToDomain(e)
		users = append(users, u)
	}
	return users, nil
}
