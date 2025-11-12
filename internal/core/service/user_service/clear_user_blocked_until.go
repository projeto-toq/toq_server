package userservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ClearUserBlockedUntil removes temporal blocking from users table
//
// This method clears users.blocked_until (sets to NULL) without modifying user_roles.status,
// restoring authentication access while preserving validation state.
//
// Parameters:
//   - ctx: Context with logger and tracing
//   - tx: Database transaction (REQUIRED for atomic operation)
//   - userID: ID of the user to unblock
//
// Returns:
//   - error: Repository errors (idempotent - returns sql.ErrNoRows if already unblocked)
//
// Business Rules:
//   - Only clears temporal block (users.blocked_until = NULL)
//   - Does NOT modify user_roles table (preserves status)
//   - Idempotent operation (no-op if already unblocked)
//   - Can be called on successful login or by worker for expired blocks
//
// Database Operations:
//   - UPDATE users SET blocked_until = NULL WHERE id = ? AND blocked_until IS NOT NULL AND deleted = 0
//
// Side Effects:
//   - Restores user authentication access
//   - Logs INFO when successfully clears block
//   - Logs DEBUG for idempotent calls (already unblocked)
//
// Error Handling:
//   - Returns sql.ErrNoRows if user already unblocked (idempotent semantics)
//   - Infrastructure errors logged and returned
//   - Caller should ignore sql.ErrNoRows for idempotent behavior
//
// Example:
//
//	err := us.ClearUserBlockedUntil(ctx, tx, userID)
//	if err != nil && !errors.Is(err, sql.ErrNoRows) {
//	    logger.Warn("Failed to clear block", "user_id", userID, "error", err)
//	}
func (us *userService) ClearUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	err := us.repo.ClearUserBlockedUntil(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Idempotent: user already unblocked
			logger.Debug("user_service.clear_blocked_until_idempotent",
				"user_id", userID)
			return err // Let caller decide if this is an error
		}

		logger.Error("user_service.clear_blocked_until_failed",
			"user_id", userID,
			"error", err)
		return err
	}

	logger.Info("user_service.clear_blocked_until_success",
		"user_id", userID)
	return nil
}
