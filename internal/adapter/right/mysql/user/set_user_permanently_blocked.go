package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SetUserPermanentlyBlocked sets or clears permanent admin block for a user
//
// This function controls the permanently_blocked flag in users table.
// Used by admin endpoints to permanently block/unblock users (policy violations, fraud, etc.).
//
// NEW ARCHITECTURE:
//   - Blocking is now at USER level (not user_role level)
//   - Does NOT affect user_roles.status (preserves validation state)
//   - Blocks ALL roles of the user (not just active role)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - userID: User ID to block/unblock permanently
//   - blocked: true = block permanently, false = unblock
//
// Returns:
//   - error: sql.ErrNoRows if user not found, database errors
//
// Business Rules:
//   - Updates users.permanently_blocked column (NEW location)
//   - Permanent block has NO expiration (requires manual admin action to unblock)
//   - Does NOT modify user_roles.status (preserves pending states)
//   - When unblocking, also clears blocked_until (removes temporal block)
//
// Database Schema:
//   - Table: users
//   - Filter: id = ? AND deleted = 0
//   - Columns updated: permanently_blocked, blocked_until (when unblocking)
//
// Example:
//
//	// Admin permanently blocks user for fraud
//	err := adapter.SetUserPermanentlyBlocked(ctx, tx, userID, true)
//
//	// Admin unblocks user after review
//	err := adapter.SetUserPermanentlyBlocked(ctx, tx, userID, false)
func (ua *UserAdapter) SetUserPermanentlyBlocked(ctx context.Context, tx *sql.Tx, userID int64, blocked bool) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	var query string
	if blocked {
		// Block permanently (set flag to 1)
		query = `UPDATE users SET permanently_blocked = 1 WHERE id = ? AND deleted = 0`
	} else {
		// Unblock permanently (set flag to 0 and clear blocked_until)
		query = `UPDATE users SET permanently_blocked = 0, blocked_until = NULL WHERE id = ? AND deleted = 0`
	}

	result, execErr := ua.ExecContext(ctx, tx, "update", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.set_permanently_blocked.exec_error", "user_id", userID, "blocked", blocked, "error", execErr)
		return fmt.Errorf("set user permanently blocked: %w", execErr)
	}

	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.set_permanently_blocked.rows_affected_error", "user_id", userID, "error", raErr)
		return fmt.Errorf("set user permanently blocked rows affected: %w", raErr)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.set_permanently_blocked.success", "user_id", userID, "blocked", blocked)
	return nil
}
