package userservices

import (
	"context"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetExpiredTempBlockedUsers returns all users whose temporary block has expired
//
// This method retrieves user roles with StatusTempBlocked where blocked_until <= NOW().
// Used by TempBlockCleanerWorker to identify users that should be automatically unblocked.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//
// Returns:
//   - users: Slice of UserRoleInterface with expired blocks (empty if none)
//   - error: Infrastructure errors (database, transaction) mapped to InternalError (500)
//
// Business Rules:
//   - Returns only roles with blocked_until <= NOW() (expired blocks)
//   - Returns empty slice if no expired blocks (NOT an error)
//   - Includes user_id, role_id, status, blocked_until in result
//
// Database Operations:
//   - SELECT from user_roles - via GetExpiredTempBlockedUsers adapter
//   - Uses transaction for consistency
//
// Side Effects:
//   - None (read-only operation)
//   - Logs DEBUG when starting
//   - Logs ERROR if operations fail
//
// Error Handling:
//   - Infrastructure errors logged as ERROR
//   - Errors marked in span for distributed tracing
//   - Returns generic InternalError to prevent information disclosure
//   - Transaction automatically rolled back on error
//
// Observability:
//   - Log entry on start: "permission.user.temp_block.get_expired.start" (DEBUG)
//   - Log entry on error: "permission.user.temp_block.get_expired.db_failed" (ERROR)
//   - Span error marking for failed operations
//
// Important Notes:
//   - Creates and manages its own transaction
//   - Empty result is valid (no expired blocks to process)
//   - Worker calls this method periodically (every 5 minutes)
//
// Example (Worker):
//
//	expiredUsers, err := userService.GetExpiredTempBlockedUsers(ctx)
//	if err != nil {
//	    logger.Error("Failed to get expired blocks", "error", err)
//	    return
//	}
//	for _, userRole := range expiredUsers {
//	    // Unblock each user
//	}
func (us *userService) GetExpiredTempBlockedUsers(ctx context.Context) ([]usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.temp_block.get_expired.start")

	// Start transaction for read consistency
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user.temp_block.get_expired.tx_start_failed", "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user.temp_block.get_expired.tx_rollback_failed", "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		} else {
			if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
				logger.Error("permission.user.temp_block.get_expired.tx_commit_failed", "error", cmErr)
				utils.SetSpanError(ctx, cmErr)
			}
		}
	}()

	// Retrieve expired blocks from database
	users, err := us.repo.GetExpiredTempBlockedUsers(ctx, tx)
	if err != nil {
		logger.Error("permission.user.temp_block.get_expired.db_failed", "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return users, nil
}
