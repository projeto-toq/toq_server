package userservices

import (
	"context"
	"database/sql"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SetUserBlockedUntil sets temporal blocking expiration on users table
//
// This method updates users.blocked_until directly without touching user_roles.status,
// preserving any validation progress (e.g., StatusPendingPhone).
//
// Parameters:
//   - ctx: Context with logger and tracing
//   - tx: Database transaction (REQUIRED for atomic operation)
//   - userID: ID of the user to block temporarily
//   - blockedUntil: UTC timestamp when block expires
//
// Returns:
//   - error: Repository errors (user not found, database issues)
//
// Business Rules:
//   - Only sets temporal block (users.blocked_until column)
//   - Does NOT modify user_roles table (preserves status)
//   - Block expires automatically when timestamp passes
//   - Worker process clears expired blocks periodically
//
// Database Operations:
//   - UPDATE users SET blocked_until = ? WHERE id = ? AND deleted = 0
//
// Side Effects:
//   - User cannot authenticate until blockedUntil expires
//   - Logs ERROR if operation fails
//
// Error Handling:
//   - Returns sql.ErrNoRows if user not found
//   - Infrastructure errors logged and returned
//   - Caller responsible for error mapping to HTTP status
//
// Example:
//
//	blockedUntil := time.Now().UTC().Add(15 * time.Minute)
//	err := us.SetUserBlockedUntil(ctx, tx, userID, blockedUntil)
//	if err != nil {
//	    return utils.InternalError("Failed to process security measures")
//	}
func (us *userService) SetUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	err := us.repo.SetUserBlockedUntil(ctx, tx, userID, blockedUntil)
	if err != nil {
		logger.Error("user_service.set_blocked_until_failed",
			"user_id", userID,
			"blocked_until", blockedUntil,
			"error", err)
		return err
	}

	return nil
}
