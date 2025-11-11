package userservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UnblockUser unblocks a temporarily blocked user and resets failed signin attempts
//
// This method orchestrates the complete unblock flow:
//  1. Restore user_roles.status to StatusActive
//  2. Clear user_roles.blocked_until timestamp
//  3. Delete temp_wrong_signin tracking record
//  4. Clear user permissions cache
//
// Used by TempBlockCleanerWorker for automatic unblocks when blocked_until expires.
// Can also be called by admin endpoints for manual unblocks.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for atomic operations)
//   - userID: ID of the user to unblock
//
// Returns:
//   - error: Infrastructure errors (database, cache) mapped to InternalError (500)
//
// Business Rules:
//   - Only unblocks users with StatusTempBlocked (not permanent blocks)
//   - Clears failed signin attempt counter (fresh start)
//   - Cache invalidation ensures permissions refresh on next request
//   - Operation is atomic (all DB changes committed together)
//
// Database Operations:
//   - UPDATE user_roles (status + blocked_until) - via UnblockUser adapter
//   - DELETE temp_wrong_signin - via ResetUserWrongSigninAttempts adapter
//
// Side Effects:
//   - Modifies user_roles table (status restored, blocked_until cleared)
//   - Modifies temp_wrong_signin table (record deleted)
//   - Invalidates Redis cache (permissions cleared)
//   - Logs INFO when unblock succeeds
//   - Logs ERROR if operations fail
//
// Error Handling:
//   - Infrastructure errors logged as ERROR with context
//   - Errors marked in span for distributed tracing
//   - Returns generic InternalError to prevent information disclosure
//   - Transaction rollback handled by caller
//
// Observability:
//   - Log entry on error: "permission.user.unblock.db_failed"
//   - Log entry on cache error: "permission.user.unblock.cache_clear_failed" (WARN)
//   - Log entry on success: "permission.user.unblocked" (INFO)
//   - Span error marking for failed operations
//
// Important Notes:
//   - MUST be called within a transaction (tx parameter)
//   - Caller MUST manage transaction lifecycle (commit/rollback)
//   - Cache clearing failure is non-fatal (logged as WARN)
//   - Used ONLY by worker and future admin endpoints
//
// Example (Worker):
//
//	tx, _ := globalService.StartTransaction(ctx)
//	defer rollbackOnError(tx, err)
//
//	err = userService.UnblockUser(ctx, tx, userID)
//	if err != nil {
//	    return  // Error already logged
//	}
//
//	_ = globalService.CommitTransaction(ctx, tx)
func (us *userService) UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.unblock.start", "user_id", userID)

	// Step 1: Restore user role status and clear blocked_until
	err = us.repo.UnblockUser(ctx, tx, userID)
	if err != nil {
		logger.Error("permission.user.unblock.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("Failed to unblock user")
	}

	// Step 2: Reset failed signin attempts counter (now transactional)
	err = us.repo.ResetUserWrongSigninAttempts(ctx, tx, userID)
	if err != nil {
		logger.Error("permission.user.unblock.reset_attempts_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("Failed to reset signin attempts")
	}

	// Step 3: Clear user permissions cache to force refresh
	// Note: Cache clearing failure is non-fatal (logged as WARN, not ERROR)
	if errCache := us.permissionService.ClearUserPermissionsCache(ctx, userID); errCache != nil {
		logger.Warn("permission.user.unblock.cache_clear_failed", "user_id", userID, "error", errCache)
		// Don't return error here as the unblocking was successful
	}

	// Success: log INFO (domain event)
	logger.Info("permission.user.unblocked", "user_id", userID)
	return nil
}
