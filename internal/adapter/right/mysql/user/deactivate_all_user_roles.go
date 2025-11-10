package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeactivateAllUserRoles sets all user roles to inactive status
//
// This function updates the is_active flag to 0 for ALL role assignments
// belonging to a specific user. This is typically used when suspending a user
// account or during role reassignment workflows.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - role deactivation must be transactional)
//   - userID: ID of the user whose roles will be deactivated
//
// Returns:
//   - error: sql.ErrNoRows if user has no roles, database errors otherwise
//
// Business Rules:
//   - Returns sql.ErrNoRows if user has NO role assignments (mapped to 404 by service)
//   - Deactivates ALL roles (active and already inactive)
//   - Does NOT delete role records (only sets is_active = 0)
//
// Database Constraints:
//   - FK: user_id REFERENCES users(id) ON DELETE CASCADE
//
// Usage Example:
//
//	err := adapter.DeactivateAllUserRoles(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    return derrors.NotFound("User has no roles")
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) DeactivateAllUserRoles(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update all user roles to inactive (bulk deactivation)
	query := `
		UPDATE user_roles
		SET is_active = 0
		WHERE user_id = ?
	`

	// Execute update using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "update", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.deactivate_all_user_roles.exec_error", "error", execErr)
		return fmt.Errorf("execute deactivate all user roles: %w", execErr)
	}

	// Check if user has any role assignments
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.deactivate_all_user_roles.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("rows affected deactivate all user roles: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user has no roles (standard repository pattern)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.deactivate_all_user_roles.success", "rows_affected", rowsAffected)
	return nil
}
