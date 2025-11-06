package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserRoleStatusByUserID updates the status of the active user role
//
// This function updates the status field in user_roles table for the currently
// active role of a user. Used to change role status (pending, active, suspended, etc.)
// without affecting role assignment.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - userID: User ID whose active role status should be updated
//   - status: New status value (0=pending, 1=active, 2=suspended, etc.)
//
// Returns:
//   - error: sql.ErrNoRows if no active role found, database errors
//
// Business Rules:
//   - Updates ONLY the active role (WHERE is_active = 1)
//   - User must have exactly one active role
//   - Status values defined by permission model enum
//   - Does NOT validate status value (service layer responsibility)
//
// Database Schema:
//   - Table: user_roles
//   - Filter: is_active = 1 (ensures only active role updated)
//   - Columns: id, user_id, role_id, is_active, status, expires_at, blocked_until
//
// Edge Cases:
//   - Returns sql.ErrNoRows if user has no active role (service maps to 404)
//   - User deleted: May still have active role (consider adding deleted check)
//   - Invalid status value: Accepted by DB (validation in service layer)
//
// Performance:
//   - Single-row UPDATE using indexed is_active column
//   - Composite index on (user_id, is_active) recommended
//
// Important Notes:
//   - Does NOT deactivate other roles (assumes only 1 active role exists)
//   - Does NOT validate role permissions after status change
//   - Transaction parameter intentionally nil (standalone operation)
//
// Example:
//
//	// Suspend user's active role
//	err := adapter.UpdateUserRoleStatusByUserID(ctx, userID, 2)
//	if err == sql.ErrNoRows {
//	    // User has no active role
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) UpdateUserRoleStatusByUserID(ctx context.Context, userID int64, status int) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update status of active role for this user
	// Note: WHERE is_active = 1 ensures only active role is updated
	query := `UPDATE user_roles SET status = ? WHERE user_id = ? AND is_active = 1`

	// Execute update using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, nil, "update", query, status, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_role_status.exec_error", "user_id", userID, "status", status, "error", execErr)
		return fmt.Errorf("update user role status by user: %w", execErr)
	}

	// Check if active role exists and was updated
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_role_status.rows_affected_error", "user_id", userID, "status", status, "error", rowsErr)
		return fmt.Errorf("update user role status rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user has no active role (service layer maps to 404)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
