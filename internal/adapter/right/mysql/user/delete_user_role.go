package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteUserRole removes a specific user-role assignment by its unique ID
//
// This function permanently deletes a user_roles record. This is typically
// called when explicitly removing a role from a user or during cleanup of
// expired role assignments.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - role deletion must be transactional)
//   - userRoleID: Primary key of the user_roles record to delete
//
// Returns:
//   - error: sql.ErrNoRows if user role not found, database errors otherwise
//
// Business Rules:
//   - Returns sql.ErrNoRows if user-role record does not exist
//   - Service layer maps sql.ErrNoRows to domain error (404 Not Found)
//   - Permanently removes the role assignment (not just deactivation)
//
// Database Constraints:
//   - FK: user_id REFERENCES users(id) ON DELETE CASCADE
//   - FK: role_id REFERENCES roles(id) ON DELETE CASCADE
//
// Usage Example:
//
//	err := adapter.DeleteUserRole(ctx, tx, userRoleID)
//	if err == sql.ErrNoRows {
//	    return derrors.NotFound("User role not found")
//	}
func (ua *UserAdapter) DeleteUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete user-role record by primary key
	query := `
		DELETE FROM user_roles 
		WHERE id = ?
	`

	// Execute deletion using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userRoleID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_user_role.exec_error", "error", execErr)
		return fmt.Errorf("delete user role: %w", execErr)
	}

	// Check if user-role record was found and deleted
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_user_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user role delete rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user role not found (standard repository pattern)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.delete_user_role.success")
	return nil
}
