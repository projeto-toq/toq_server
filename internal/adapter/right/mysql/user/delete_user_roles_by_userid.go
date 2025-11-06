package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteUserRolesByUserID removes ALL role assignments for a specific user
//
// This function performs a cascading deletion of all user_roles records
// for a given user. This is typically used when deleting a user account
// or during role reassignment workflows.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - deletion must be transactional)
//   - userID: ID of the user whose roles will be deleted
//
// Returns:
//   - deleted: Number of role records deleted
//   - error: sql.ErrNoRows if user has no roles, database errors otherwise
//
// Business Rules:
//   - Returns sql.ErrNoRows if user has NO roles (not an error condition)
//   - Service layer maps sql.ErrNoRows to appropriate domain response
//   - Deletes ALL roles (active and inactive) for the user
//
// Database Constraints:
//   - FK: user_id REFERENCES users(id) ON DELETE CASCADE
//
// Usage Example:
//
//	deleted, err := adapter.DeleteUserRolesByUserID(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    // User had no roles (expected scenario)
//	} else if err != nil {
//	    // Infrastructure error
//	}
//	// deleted contains count of removed roles
func (ua *UserAdapter) DeleteUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete all role assignments for the user
	query := `DELETE FROM user_roles WHERE user_id = ?;`

	// Execute deletion using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_user_roles.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete user_roles by user_id: %w", execErr)
	}

	// Check how many role assignments were deleted
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_user_roles.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete user_roles rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user had no roles (standard repository pattern)
	// ✅ CHANGED: Previously returned errors.New(), now returns sql.ErrNoRows
	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return rowsAffected, nil
}
