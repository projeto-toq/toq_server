package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ActivateUserRole sets a specific user role to active status
//
// This function updates the is_active flag to 1 for a specific user-role relationship.
// It is typically called when a user switches between multiple roles or when reactivating
// a previously deactivated role.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - role activation must be transactional)
//   - userID: ID of the user whose role is being activated
//   - roleID: ID of the role to activate
//
// Returns:
//   - error: sql.ErrNoRows if user-role relationship not found, database errors otherwise
//
// Business Rules:
//   - Returns sql.ErrNoRows if no matching user-role record exists (mapped to 404 by service)
//   - Does NOT check if role is already active (idempotent operation)
//   - Service layer validates that user owns this role before calling
//
// Database Constraints:
//   - FK: user_id REFERENCES users(id) ON DELETE CASCADE
//   - FK: role_id REFERENCES roles(id) ON DELETE CASCADE
//
// Usage Example:
//
//	err := adapter.ActivateUserRole(ctx, tx, userID, roleID)
//	if err == sql.ErrNoRows {
//	    return derrors.NotFound("User role not found")
//	}
func (ua *UserAdapter) ActivateUserRole(ctx context.Context, tx *sql.Tx, userID, roleID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update user_role to active status (idempotent operation)
	query := `
		UPDATE user_roles
		SET is_active = 1
		WHERE user_id = ? AND role_id = ?
	`

	// Execute update using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "update", query, userID, roleID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.activate_user_role.exec_error", "error", execErr)
		return fmt.Errorf("execute activate user role: %w", execErr)
	}

	// Check if user-role relationship exists
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.activate_user_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("rows affected activate user role: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user-role relationship not found (standard repository pattern)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.activate_user_role.success", "rows_affected", rowsAffected)
	return nil
}
