package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SetUserBlockedUntil sets temporary block expiration timestamp for a user
//
// This function blocks a user until the specified timestamp by setting users.blocked_until.
// Used for temporary security restrictions (failed signin attempts, suspicious activity).
//
// NEW ARCHITECTURE:
//   - Blocking is now at USER level (not user_role level)
//   - Does NOT affect user_roles.status (preserves validation state)
//   - Blocks ALL roles of the user (not just active role)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - userID: User ID to block temporarily
//   - blockedUntil: Timestamp when block should automatically expire
//
// Returns:
//   - error: sql.ErrNoRows if user not found, database errors
//
// Business Rules:
//   - Updates users.blocked_until column (NEW location)
//   - User is blocked while blocked_until > NOW()
//   - Worker automatically unblocks when timestamp expires
//   - Does NOT modify user_roles.status (preserves pending states)
//
// Database Schema:
//   - Table: users
//   - Filter: id = ? AND deleted = 0
//   - Column updated: blocked_until
//
// Example:
//
//	// Block user for 15 minutes after 3 failed signin attempts
//	expiresAt := time.Now().UTC().Add(15 * time.Minute)
//	err := adapter.SetUserBlockedUntil(ctx, tx, userID, expiresAt)
func (ua *UserAdapter) SetUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update users.blocked_until (NEW: moved from user_roles)
	query := `UPDATE users SET blocked_until = ? WHERE id = ? AND deleted = 0`

	result, execErr := ua.ExecContext(ctx, tx, "update", query, blockedUntil, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.set_blocked_until.exec_error", "user_id", userID, "error", execErr)
		return fmt.Errorf("set user blocked until: %w", execErr)
	}

	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.set_blocked_until.rows_affected_error", "user_id", userID, "error", raErr)
		return fmt.Errorf("set user blocked until rows affected: %w", raErr)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.set_blocked_until.success", "user_id", userID, "blocked_until", blockedUntil)
	return nil
}
