package userservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SetUserPermanentlyBlocked sets or clears permanent block flag on users table
//
// This method updates users.permanently_blocked flag for admin-initiated blocking.
// When unblocking, it also clears any temporal block to fully restore access.
//
// Parameters:
//   - ctx: Context with logger and tracing
//   - tx: Database transaction (REQUIRED for atomic operation)
//   - userID: ID of the user to block/unblock permanently
//   - blocked: true to block, false to unblock
//
// Returns:
//   - error: Repository errors (user not found, database issues)
//
// Business Rules:
//   - Sets users.permanently_blocked = 1 (block) or 0 (unblock)
//   - When unblocking (blocked=false), also clears users.blocked_until
//   - Does NOT modify user_roles table (preserves status)
//   - Intended for admin endpoints (not yet implemented)
//
// Database Operations:
//   - Block: UPDATE users SET permanently_blocked = 1 WHERE id = ? AND deleted = 0
//   - Unblock: UPDATE users SET permanently_blocked = 0, blocked_until = NULL WHERE id = ? AND deleted = 0
//
// Side Effects:
//   - Blocked users cannot authenticate (IsBlocked() returns true)
//   - Unblocking fully restores authentication access
//   - Logs INFO when operation succeeds
//
// Error Handling:
//   - Returns sql.ErrNoRows if user not found
//   - Infrastructure errors logged and returned
//   - Caller responsible for error mapping to HTTP status
//
// Example:
//
//	// Admin blocks user permanently
//	err := us.SetUserPermanentlyBlocked(ctx, tx, userID, true)
//	if err != nil {
//	    return utils.InternalError("Failed to block user")
//	}
func (us *userService) SetUserPermanentlyBlocked(ctx context.Context, tx *sql.Tx, userID int64, blocked bool) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	err := us.repo.SetUserPermanentlyBlocked(ctx, tx, userID, blocked)
	if err != nil {
		logger.Error("user_service.set_permanently_blocked_failed",
			"user_id", userID,
			"blocked", blocked,
			"error", err)
		return err
	}

	action := "unblocked"
	if blocked {
		action = "blocked"
	}
	logger.Info("user_service.set_permanently_blocked_success",
		"user_id", userID,
		"action", action)
	return nil
}
