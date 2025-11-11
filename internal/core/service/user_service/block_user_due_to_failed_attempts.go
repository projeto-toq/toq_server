package userservices

import (
	"context"
	"database/sql"
	"time"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// blockUserDueToFailedAttempts blocks a user account temporarily after excessive signin failures
//
// This method sets user's active role to StatusTempBlocked with expiration timestamp
// (via user_roles table). The lockout moment is tracked by user_roles.blocked_until.
//
// The function is called ONLY when failed attempt counter reaches MaxWrongSigninAttempts threshold.
// Blocking is temporary and automatically expires after TempBlockDuration period.
//
// Parameters:
//   - ctx: Context for logging (must contain logger from parent)
//   - tx: Database transaction (REQUIRED for atomic blocking operations)
//   - userID: ID of the user to block
//
// Returns:
//   - error: Infrastructure errors (database, transaction) mapped to InternalError (500)
//
// Business Rules:
//   - Block duration: TempBlockDuration (15 minutes from Now)
//   - Updates user_roles.status to StatusTempBlocked
//   - Sets user_roles.blocked_until with expiration timestamp
//   - Worker process automatically unblocks when blocked_until expires
//
// Database Operations:
//   - UPDATE user_roles (status + blocked_until columns) - via BlockUserTemporarily
//
// Side Effects:
//   - Modifies user_roles table (sets temp block status + expiration)
//   - User cannot authenticate until blocked_until expires
//   - Logs ERROR if operations fail
//
// Error Handling:
//   - Infrastructure errors logged as ERROR with context (user_id, timestamps)
//   - Errors marked in span for distributed tracing
//   - Returns generic InternalError to prevent information disclosure
//
// Observability:
//   - Errors logged with "auth.signin.block_user_failed" key
//   - Span error marking for failed operations
//
// Important Notes:
//   - Caller (processFailedSigninAttempt) logs WARN when blocking succeeds
//   - This function only handles blocking logic, not logging success
//   - Transaction commit/rollback handled by top-level signIn function
//
// Example:
//
//	// Called when threshold reached (3 failed attempts)
//	err := us.blockUserDueToFailedAttempts(ctx, tx, userID)
//	if err != nil {
//	    return err  // Error already logged and mapped
//	}
func (us *userService) blockUserDueToFailedAttempts(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Reuse context from parent (no new tracer for private methods)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Calculate block expiration time
	blockedUntil := time.Now().UTC().Add(usermodel.TempBlockDuration)

	// Block user's active role temporarily
	err := us.repo.BlockUserTemporarily(ctx, tx, userID, blockedUntil, "Too many failed signin attempts")
	if err != nil {
		logger.Error("auth.signin.block_user_failed",
			"user_id", userID,
			"blocked_until", blockedUntil,
			"error", err)
		return utils.InternalError("Failed to process security measures")
	}

	return nil
}
