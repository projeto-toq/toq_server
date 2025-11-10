package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserRoleByUserIDAndRoleID retrieves a specific user_role by composite key (user_id + role_id)
//
// This function searches for a user's role assignment by the combination of user ID and role ID,
// regardless of whether the assignment is active or inactive. Used for checking if a user has
// ever been assigned a specific role, or for retrieving historical assignments.
//
// Query Logic:
//   - Searches user_roles table by composite key (user_id + role_id)
//   - Does NOT filter by is_active (returns active or inactive assignments)
//   - Does NOT filter by expires_at (returns expired assignments)
//   - LIMIT 1 ensures single result (composite key should be unique)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - userID: User's unique identifier
//   - roleID: Role's unique identifier
//
// Returns:
//   - userRole: UserRoleInterface with assignment details, or nil if not found
//   - error: Database errors (returns nil instead of sql.ErrNoRows when not found)
//
// Business Rules:
//   - Composite key (user_id + role_id) should be unique per user
//   - Returns assignment regardless of status (active/inactive/expired)
//   - Does NOT populate Role aggregate (only UserRole fields)
//   - Service layer decides if inactive/expired assignment is acceptable
//
// Edge Cases:
//   - User never had this role: Returns nil userRole, nil error
//   - User had role but it's now inactive: Returns UserRole with is_active = 0
//   - User had role but it expired: Returns UserRole with expires_at < NOW()
//   - Multiple assignments (data integrity issue): Returns first by implicit order
//
// Performance:
//   - Uses composite index on (user_id, role_id) for fast lookup
//   - Single row maximum (no pagination needed)
//
// Use Cases:
//   - Checking if user has specific role assignment (for activation/reactivation)
//   - Retrieving role assignment to update status
//   - Audit queries: "Has this user ever been assigned this role?"
//
// Difference from GetActiveUserRoleByUserID:
//   - GetActiveUserRoleByUserID: Finds ONE active role, any role_id
//   - GetUserRoleByUserIDAndRoleID: Finds specific role_id, any status
//
// Example:
//
//	userRole, err := adapter.GetUserRoleByUserIDAndRoleID(ctx, tx, 123, 5)
//	if err != nil {
//	    // Handle infrastructure error
//	}
//	if userRole == nil {
//	    // User has never been assigned role ID 5
//	}
//	if !userRole.IsActive() {
//	    // Assignment exists but is inactive
//	}
func (ua *UserAdapter) GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query user_role by composite key (user_id + role_id)
	// Note: No filter by is_active - returns active or inactive assignments
	// Note: No filter by expires_at - returns expired assignments
	query := `
		SELECT id, user_id, role_id, is_active, status, expires_at
		FROM user_roles 
		WHERE user_id = ? AND role_id = ?
		LIMIT 1
	`

	// Declare strongly-typed variables for scanning
	var (
		id          int64
		uid         int64
		roleIDOut   int64
		isActiveInt int64
		status      int64
		expiresAt   sql.NullTime
	)

	// Execute query using instrumented adapter
	row := ua.QueryRowContext(ctx, tx, "select", query, userID, roleID)
	err = row.Scan(
		&id, &uid, &roleIDOut, &isActiveInt, &status, &expiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// No assignment found is not an error condition - return nil instead
			logger.Debug("mysql.user.get_user_role_by_user_id_and_role_id.not_found",
				"user_id", userID, "role_id", roleID)
			return nil, nil
		}
		// Infrastructure error: log and mark span
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_role_by_user_id_and_role_id.scan_error",
			"user_id", userID, "role_id", roleID, "error", err)
		return nil, fmt.Errorf("get user role by user id and role id scan: %w", err)
	}

	// Build strongly-typed entity from scanned values
	entity := &userentity.UserRoleEntity{
		ID:       uint32(id),
		UserID:   uint32(userID),
		RoleID:   uint32(roleIDOut),
		IsActive: isActiveInt == 1,
		Status:   int8(status),
	}
	if expiresAt.Valid {
		entity.ExpiresAt = sql.NullTime{
			Time:  expiresAt.Time,
			Valid: true,
		}
	}

	// Convert entity to domain model using type-safe converter
	userRole, convertErr := userconverters.UserRoleEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.user.get_user_role_by_user_id_and_role_id.convert_error",
			"user_id", userID, "role_id", roleID, "error", convertErr)
		return nil, fmt.Errorf("convert user role entity to domain: %w", convertErr)
	}

	logger.Debug("mysql.user.get_user_role_by_user_id_and_role_id.success",
		"user_id", userID, "role_id", roleID, "is_active", userRole.GetIsActive())
	return userRole, nil
}
