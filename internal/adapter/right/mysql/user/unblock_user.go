package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UnblockUser unblocks a temporarily blocked user by setting their status back to active and clearing blocked_until
//
// This function removes temporary block restrictions from a user's active role.
// Used for manual unblocks by admin or automatic unblocks by worker when blocked_until expires.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - userID: User ID to unblock
//
// Returns:
//   - error: sql.ErrNoRows if no temp blocked role found, database errors
//
// Business Rules:
//   - Updates ONLY active role with StatusTempBlocked (WHERE status = ? AND is_active = 1)
//   - Sets status to StatusActive (restores authentication access)
//   - Clears blocked_until timestamp (NULL)
//   - Does NOT unblock permanently blocked users (different status)
//
// Database Schema:
//   - Table: user_roles
//   - Filter: is_active = 1 AND status = StatusTempBlocked
//   - Columns updated: status, blocked_until
//
// Edge Cases:
//   - Returns sql.ErrNoRows if user not temp blocked (service maps to 404)
//   - User permanently blocked: Not affected by this method
//   - User already active: No rows updated (returns sql.ErrNoRows)
//
// Performance:
//   - Single-row UPDATE using indexed columns
//   - Composite index on (user_id, is_active, status) recommended
//
// Important Notes:
//   - Only unblocks StatusTempBlocked users (not permanent blocks)
//   - Service layer should log unblock event in audit table
//   - Does NOT send notification to user (service layer responsibility)
//
// Example:
//
//	// Manually unblock user (admin action)
//	err := adapter.UnblockUser(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    // User has no active temp blocked role
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Restore active status and clear blocked_until for temp blocked users
	// Note: WHERE status = StatusTempBlocked ensures only temp blocks are cleared
	// Note: Does NOT unblock permanently blocked users (different status value)
	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = NULL
		WHERE user_id = ? AND status = ? AND is_active = 1
	`

	// Execute update using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		globalmodel.StatusActive,
		userID,
		globalmodel.StatusTempBlocked,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.unblock_user.exec_error", "user_id", userID, "error", execErr)
		return fmt.Errorf("unblock user exec: %w", execErr)
	}

	// Check if temp blocked role exists and was unblocked
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.unblock_user.rows_affected_error", "user_id", userID, "error", rowsErr)
		return fmt.Errorf("unblock user rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user not temp blocked (service layer maps to 404)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.unblock_user.success", "user_id", userID)
	return nil
}
