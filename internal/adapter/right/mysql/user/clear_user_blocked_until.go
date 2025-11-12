package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ClearUserBlockedUntil clears temporary block for a user
//
// This function removes temporary block by setting users.blocked_until = NULL.
// Used by worker when block expires, or by signin flow on successful authentication.
//
// NEW ARCHITECTURE:
//   - Blocking is now at USER level (not user_role level)
//   - Does NOT modify user_roles.status (preserves validation state)
//   - Idempotent: safe to call multiple times (no error if already unblocked)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - userID: User ID to unblock
//
// Returns:
//   - error: sql.ErrNoRows if user has no blocked_until set, database errors
//
// Business Rules:
//   - Updates users.blocked_until = NULL (NEW location)
//   - Idempotent: returns sql.ErrNoRows if already unblocked
//   - Service should treat sql.ErrNoRows as success (idempotent operation)
//
// Database Schema:
//   - Table: users
//   - Filter: id = ? AND deleted = 0 AND blocked_until IS NOT NULL
//   - Column updated: blocked_until
//
// Example:
//
//	// Worker unblocks user after blocked_until expires
//	err := adapter.ClearUserBlockedUntil(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    // User already unblocked (idempotent)
//	} else if err != nil {
//	    // Infrastructure error
//	}
func (ua *UserAdapter) ClearUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Clear users.blocked_until (idempotent: only updates if blocked_until IS NOT NULL)
	query := `
		UPDATE users 
		SET blocked_until = NULL 
		WHERE id = ? 
		  AND deleted = 0 
		  AND blocked_until IS NOT NULL
	`

	result, execErr := ua.ExecContext(ctx, tx, "update", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.clear_blocked_until.exec_error", "user_id", userID, "error", execErr)
		return fmt.Errorf("clear user blocked until: %w", execErr)
	}

	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.clear_blocked_until.rows_affected_error", "user_id", userID, "error", raErr)
		return fmt.Errorf("clear user blocked until rows affected: %w", raErr)
	}

	// Return sql.ErrNoRows if user has no blocked_until (idempotent)
	if rowsAffected == 0 {
		logger.Debug("mysql.user.clear_blocked_until.no_rows", "user_id", userID)
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.clear_blocked_until.success", "user_id", userID)
	return nil
}
