package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserRoleStatus updates the active user role status for a specific role using a transaction
//
// This function updates the status field in user_roles table for the currently active role
// matching a specific role slug. Unlike UpdateUserRoleStatusByUserID (which updates by user_id only),
// this method requires BOTH user_id AND role_slug, allowing precise control when a user has
// multiple role assignments.
//
// Use Cases:
//   - Update status of specific role when user has multiple roles
//   - Change role status based on role type (e.g., suspend only realtor role, keep owner active)
//   - Atomic updates within larger transactions
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for ACID guarantees)
//   - userID: User ID whose role status should be updated
//   - role: Role slug identifier (e.g., "owner", "realtor", "agency")
//   - status: New status value (StatusPending, StatusActive, StatusSuspended, etc.)
//
// Returns:
//   - error: sql.ErrNoRows if no active role found matching criteria, database errors
//
// Business Rules:
//   - Updates ONLY the active role (WHERE ur.is_active = 1)
//   - Role must match provided slug exactly
//   - Status values defined by globalmodel.UserRoleStatus enum
//   - Does NOT validate status transition logic (service layer responsibility)
//
// Query Structure:
//   - JOIN user_roles with roles to match slug
//   - Filter: ur.user_id = ? AND ur.is_active = 1 AND r.slug = ?
//   - Updates: ur.status only (does not affect is_active or other fields)
//
// Database Schema:
//   - Table: user_roles (ur) JOIN roles (r)
//   - Filter: is_active = 1 ensures only active role updated
//   - Index: Composite index on (user_id, is_active, role_id) recommended
//
// Edge Cases:
//   - Returns sql.ErrNoRows if user has no active role with this slug
//   - User deleted: May still have active role (consider adding deleted check in service)
//   - Invalid role slug: No rows updated (returns sql.ErrNoRows)
//   - Multiple active roles with same slug: Should never happen (data integrity issue)
//
// Performance:
//   - Single-row UPDATE using indexed columns
//   - JOIN with roles is efficient (role_id indexed)
//
// Comparison with UpdateUserRoleStatusByUserID:
//   - UpdateUserRoleStatusByUserID: Updates active role for user (any role)
//   - UpdateUserRoleStatus: Updates active role matching SPECIFIC slug
//
// Important Notes:
//   - Requires transaction for consistency with related operations
//   - Does NOT deactivate other roles
//   - Does NOT validate role permissions after status change
//
// Example:
//
//	// Suspend user's realtor role (keep other roles active)
//	err := adapter.UpdateUserRoleStatus(ctx, tx, userID, "realtor", globalmodel.StatusSuspended)
//	if err == sql.ErrNoRows {
//	    // User has no active realtor role
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) UpdateUserRoleStatus(ctx context.Context, tx *sql.Tx, userID int64, role permissionmodel.RoleSlug, status globalmodel.UserRoleStatus) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update status of active role matching specific slug
	// JOIN with roles table allows filtering by role slug instead of role ID
	// Note: WHERE ur.is_active = 1 ensures only active role is updated
	// Note: r.slug = ? prevents accidental updates to wrong role type
	const query = `
		UPDATE user_roles ur
		JOIN roles r ON r.id = ur.role_id
		SET ur.status = ?
		WHERE ur.user_id = ? AND ur.is_active = 1 AND r.slug = ?`

	// Execute update using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "update", query, int(status), userID, role)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_role_status_tx.update_error",
			"user_id", userID, "role", role, "status", status, "error", execErr)
		return fmt.Errorf("update user role status: %w", execErr)
	}

	// Check if active role exists and was updated
	if result != nil {
		rowsAffected, rowsErr := result.RowsAffected()
		if rowsErr == nil {
			// No rows updated indicates absence of active role with provided slug
			if rowsAffected == 0 {
				errNoRows := sql.ErrNoRows
				utils.SetSpanError(ctx, errNoRows)
				logger.Error("mysql.user.update_user_role_status_tx.no_rows",
					"user_id", userID, "role", role, "status", status, "error", errNoRows)
				return errNoRows
			}
		} else {
			// Log warning but don't fail operation (rows affected check is optional)
			logger.Warn("mysql.user.update_user_role_status_tx.rows_affected_warning",
				"user_id", userID, "role", role, "status", status, "error", rowsErr)
		}
	}

	logger.Debug("mysql.user.update_user_role_status_tx.success",
		"user_id", userID, "role", role, "status", status)
	return nil
}
