package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserRolesByUserID retrieves all role assignments for a specific user
//
// This function performs a JOIN between user_roles and roles tables to fetch all role
// assignments for a user, regardless of their status (active/inactive/expired). The result
// includes both the assignment details (user_roles) and the role definition (roles).
//
// Query Strategy:
//   - JOINs user_roles with roles table
//   - Does NOT filter by is_active (returns all assignments)
//   - Does NOT filter by expires_at (returns expired assignments)
//   - Orders by ur.id ASC (chronological order of assignments)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - userID: User's unique identifier
//
// Returns:
//   - userRoles: Slice of UserRoleInterface with Role populated for each
//   - error: Database errors (returns empty slice instead of sql.ErrNoRows)
//
// Business Rules:
//   - Returns ALL assignments (active, inactive, expired)
//   - Each UserRoleInterface has its Role aggregate populated
//   - Order is chronological (oldest assignment first)
//   - Service layer filters by status if needed
//
// Edge Cases:
//   - User has no role assignments: Returns empty slice, no error
//   - User has multiple roles: Returns all in chronological order
//   - Some assignments expired: All returned, service decides handling
//
// Performance:
//   - Single query with INNER JOIN (efficient)
//   - Uses index on user_roles.user_id
//   - Returns all assignments in one round-trip
//
// Use Cases:
//   - User profile: Display all historical role assignments
//   - Admin panel: Manage user's role assignments
//   - Audit: Track role assignment history
//   - Authorization: Check all user's roles (including inactive for special cases)
//
// Difference from GetActiveUserRoleByUserID:
//   - GetActiveUserRoleByUserID: Returns ONE active role
//   - GetUserRolesByUserID: Returns ALL assignments (active + inactive + expired)
//
// Example:
//
//	userRoles, err := adapter.GetUserRolesByUserID(ctx, tx, 123)
//	if err != nil {
//	    // Handle infrastructure error
//	}
//	if len(userRoles) == 0 {
//	    // User has no role assignments
//	}
//	for _, ur := range userRoles {
//	    // Process each assignment with ur.GetRole() populated
//	}
func (ua *UserAdapter) GetUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query all user_roles with JOIN to roles for complete data
	// Note: INNER JOIN excludes orphaned user_roles (role deleted from roles table)
	// Note: No filter by is_active - returns all assignments
	query := `
		SELECT 
			ur.id,
			ur.user_id,
			ur.role_id,
			ur.is_active,
			ur.status,
			ur.expires_at,
			r.id,
			r.slug,
			r.name,
			r.description,
			r.is_system_role,
			r.is_active
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ?
		ORDER BY ur.id ASC
	`

	// Execute query using instrumented adapter
	rows, readErr := ua.QueryContext(ctx, tx, "select", query, userID)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.user.get_user_roles_by_user_id.query_error", "user_id", userID, "error", readErr)
		return nil, fmt.Errorf("get user roles by user id query: %w", readErr)
	}
	defer rows.Close()

	// Scan rows using strongly-typed helper function
	userRoleEntities, roleEntities, rowsErr := scanUserRoleWithRoleEntities(rows)
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.get_user_roles_by_user_id.scan_error", "user_id", userID, "error", rowsErr)
		return nil, fmt.Errorf("scan user roles with roles: %w", rowsErr)
	}

	// Convert entities to domain models (parallel arrays)
	userRoles := make([]usermodel.UserRoleInterface, 0, len(userRoleEntities))
	for i := range userRoleEntities {
		// Convert UserRoleEntity to domain
		userRole, convertErr := userconverters.UserRoleEntityToDomain(&userRoleEntities[i])
		if convertErr != nil {
			utils.SetSpanError(ctx, convertErr)
			logger.Error("mysql.user.get_user_roles_by_user_id.convert_user_role_error",
				"user_id", userID, "index", i, "error", convertErr)
			return nil, fmt.Errorf("convert user role entity to domain: %w", convertErr)
		}

		// Convert RoleEntity to domain and associate with UserRole
		if userRole != nil {
			role := permissionconverters.RoleEntityToDomain(&roleEntities[i])
			if role != nil {
				userRole.SetRole(role)
			}
			userRoles = append(userRoles, userRole)
		}
	}

	logger.Debug("mysql.user.get_user_roles_by_user_id.success", "user_id", userID, "count", len(userRoles))
	return userRoles, nil
}
