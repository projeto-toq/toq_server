package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user temporarily by updating their active user_role status and blocked_until timestamp
//
// This function sets the user's active role to StatusTempBlocked and records the unblock timestamp.
// Used for temporary account suspensions due to security issues, policy violations, or automated restrictions.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency with related operations)
//   - userID: User ID to block temporarily
//   - blockedUntil: Timestamp when block should automatically expire
//   - reason: Human-readable explanation for block (for audit/support purposes)
//
// Returns:
//   - error: sql.ErrNoRows if no active role found, database errors
//
// Business Rules:
//   - Updates ONLY the active role (WHERE is_active = 1)
//   - Sets status to StatusTempBlocked (prevents authentication)
//   - Sets blocked_until timestamp (automatic unblock time)
//   - Reason parameter is for logging/audit (not stored in DB currently)
//
// Database Schema:
//   - Table: user_roles
//   - Filter: is_active = 1 (ensures only active role blocked)
//   - Columns updated: status, blocked_until
//
// Edge Cases:
//   - Returns sql.ErrNoRows if user has no active role (service maps to 404)
//   - User already blocked: Updates blocked_until with new value (extends block)
//   - blockedUntil in past: Accepted but user will be unblocked immediately by worker
//
// Performance:
//   - Single-row UPDATE using indexed is_active column
//   - Composite index on (user_id, is_active) recommended
//
// Important Notes:
//   - Worker process monitors blocked_until and auto-unblocks expired blocks
//   - Service layer should log reason in audit table
//   - Does NOT send notification to user (service layer responsibility)
//
// Example:
//
//	// Block user for 24 hours due to failed signin attempts
//	blockedUntil := time.Now().Add(24 * time.Hour)
//	err := adapter.BlockUserTemporarily(ctx, tx, userID, blockedUntil, "Too many failed signin attempts")
//	if err == sql.ErrNoRows {
//	    // User has no active role
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) BlockUserTemporarily(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time, reason string) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update active role status to temp blocked with expiration timestamp
	// Note: WHERE is_active = 1 ensures only active role is blocked
	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = ?
		WHERE user_id = ? AND is_active = 1
	`

	// Execute update using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		globalmodel.StatusTempBlocked,
		blockedUntil,
		userID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.block_user_temporarily.exec_error",
			"user_id", userID, "blocked_until", blockedUntil, "reason", reason, "error", execErr)
		return fmt.Errorf("block user temporarily exec: %w", execErr)
	}

	// Check if active role exists and was blocked
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.block_user_temporarily.rows_affected_error",
			"user_id", userID, "error", rowsErr)
		return fmt.Errorf("block user temporarily rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user has no active role (service layer maps to 404)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.block_user_temporarily.success",
		"user_id", userID, "blocked_until", blockedUntil, "reason", reason)
	return nil
}
