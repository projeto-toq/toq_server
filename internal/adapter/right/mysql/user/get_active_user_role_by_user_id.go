package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetActiveUserRoleByUserID retrieves the user's single active role with role details
//
// This function performs a JOIN between user_roles and roles tables to fetch the user's
// currently active role assignment along with the complete role definition in a single query.
//
// Query Strategy:
//   - JOINs user_roles with roles table
//   - Filters by ur.is_active = 1 (only current active role)
//   - Filters by ur.expires_at IS NULL OR expires_at > NOW() (not expired)
//   - Filters by r.is_active = 1 (role itself is active)
//   - Orders by ur.id DESC to get latest assignment if multiple exist (data integrity issue)
//   - LIMIT 1 ensures single result
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - userID: User's unique identifier
//
// Returns:
//   - userRole: UserRoleInterface with Role populated, or nil if no active role found
//   - error: Database errors (does NOT return sql.ErrNoRows, returns nil instead)
//
// Business Rules:
//   - User has at most ONE active role at any time (enforced by service layer)
//   - Expired roles (expires_at < NOW()) are excluded
//   - Inactive roles (r.is_active = 0) are excluded
//   - Returns nil (not error) if user has no active role
//
// Edge Cases:
//   - User has no roles: Returns nil userRole, nil error
//   - User has only expired roles: Returns nil userRole, nil error
//   - User has only inactive roles: Returns nil userRole, nil error
//   - Multiple active roles (data integrity violation): Returns latest by ur.id
//
// Performance:
//   - Single query with INNER JOIN (efficient)
//   - Uses indexes: user_roles.user_id, user_roles.is_active, roles.id
//
// Use Cases:
//   - Authentication: Load user's current role and permissions
//   - Authorization: Check if user has specific active role
//   - Profile display: Show user's current role information
//
// Example:
//
//	userRole, err := adapter.GetActiveUserRoleByUserID(ctx, tx, 123)
//	if err != nil {
//	    // Handle infrastructure error
//	}
//	if userRole == nil {
//	    // User has no active role (pending activation, blocked, or no assignment)
//	}
//	// userRole.GetRole() contains complete role definition
func (ua *UserAdapter) GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query with JOIN to populate Role in UserRole in single round-trip
	// Note: INNER JOIN excludes users without roles
	// Note: Active role filters ensure only valid, non-expired assignments
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
		  AND ur.is_active = 1
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		  AND r.is_active = 1
		ORDER BY ur.id DESC
		LIMIT 1`

	// Declare strongly-typed variables for scanning
	var (
		id          int64
		uid         int64
		roleID      int64
		isActiveInt int64
		status      int64
		expiresAt   sql.NullTime

		rID          int64
		rSlug        string
		rName        string
		rDescription sql.NullString
		rIsSystemInt int64
		rIsActiveInt int64
	)

	// Execute query using instrumented adapter
	row := ua.QueryRowContext(ctx, tx, "select", query, userID)
	err = row.Scan(
		&id,
		&uid,
		&roleID,
		&isActiveInt,
		&status,
		&expiresAt,
		&rID,
		&rSlug,
		&rName,
		&rDescription,
		&rIsSystemInt,
		&rIsActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// No active role found is not an error condition - return nil instead
			logger.Debug("mysql.user.get_active_user_role_by_user_id.not_found", "user_id", userID)
			return nil, nil
		}
		// Infrastructure error: log and mark span
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_active_user_role_by_user_id.scan_error", "user_id", userID, "error", err)
		return nil, fmt.Errorf("get active user role by user id scan: %w", err)
	}

	// Build strongly-typed entities from scanned values
	userRoleEntity := &userentity.UserRoleEntity{
		ID:       uint32(id),
		UserID:   uint32(userID),
		RoleID:   uint32(roleID),
		IsActive: isActiveInt == 1,
		Status:   int8(status),
	}
	if expiresAt.Valid {
		userRoleEntity.ExpiresAt = sql.NullTime{
			Time:  expiresAt.Time,
			Valid: true,
		}
	}

	roleEntity := &permissionentities.RoleEntity{
		ID:   rID,
		Name: rName,
		Slug: rSlug,
		Description: func() string {
			if rDescription.Valid {
				return rDescription.String
			}
			return ""
		}(),
		IsSystemRole: rIsSystemInt == 1,
		IsActive:     rIsActiveInt == 1,
	}

	// Convert entities to domain models using type-safe converters
	userRole, convertErr := userconverters.UserRoleEntityToDomain(userRoleEntity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.user.get_active_user_role_by_user_id.convert_user_role_error", "error", convertErr)
		return nil, fmt.Errorf("convert active user role entity to domain: %w", convertErr)
	}

	// Associate role with user_role aggregate
	if userRole != nil {
		role := permissionconverters.RoleEntityToDomain(roleEntity)
		if role != nil {
			userRole.SetRole(role)
		}
	}

	logger.Debug("mysql.user.get_active_user_role_by_user_id.success", "user_id", userID, "role_id", roleID)
	return userRole, nil
}
