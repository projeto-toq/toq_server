package userservices

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// IsUserTempBlocked checks if a user is temporarily blocked
//
// This method verifies if the user's active role has StatusTempBlocked and if the block
// is still valid (blocked_until > NOW()). Creates its own transaction for the check.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - userID: ID of the user to check
//
// Returns:
//   - blocked: true if user is currently temp blocked, false otherwise
//   - error: Infrastructure errors (database, transaction) mapped to InternalError (500)
//
// Business Rules:
//   - Returns true ONLY if status = StatusTempBlocked AND blocked_until > NOW()
//   - Returns false if user has no active role
//   - Returns false if blocked_until has expired (block is over)
//
// Database Operations:
//   - SELECT from user_roles - via GetActiveUserRoleByUserID adapter
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
//   - Log entry on start: "permission.user.temp_block.check.start" (DEBUG)
//   - Log entry on error: "permission.user.temp_block.check.tx_start_failed" (ERROR)
//   - Span error marking for failed operations
//
// Important Notes:
//   - Creates and manages its own transaction
//   - For use when caller doesn't have a transaction
//   - Use IsUserTempBlockedWithTx when inside an existing transaction
//
// Example (Signin Flow):
//
//	isBlocked, err := userService.IsUserTempBlocked(ctx, userID)
//	if err != nil {
//	    return  // Infrastructure error
//	}
//	if isBlocked {
//	    return AuthenticationError("Invalid credentials")  // Generic error
//	}
func (us *userService) IsUserTempBlocked(ctx context.Context, userID int64) (bool, error) {
	// Initialize tracing for observability
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.temp_block.check.start", "user_id", userID)

	// Start a transaction for read operations when caller doesn't manage one
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user.temp_block.check.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user.temp_block.check.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	// Delegate to transaction-aware method
	blocked, ierr := us.IsUserTempBlockedWithTx(ctx, tx, userID)
	if ierr != nil {
		return false, ierr
	}

	// Commit transaction
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user.temp_block.check.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}

	return blocked, nil
}

// IsUserTempBlockedWithTx checks if a user is temporarily blocked using the provided transaction
//
// This method is the transactional variant of IsUserTempBlocked. Use this when you already
// have an active transaction and want to avoid creating a nested transaction.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED - must be managed by caller)
//   - userID: ID of the user to check
//
// Returns:
//   - blocked: true if user is currently temp blocked, false otherwise
//   - error: Infrastructure errors (database) mapped to InternalError (500)
//
// Business Rules:
//   - Returns true ONLY if status = StatusTempBlocked AND blocked_until > NOW()
//   - Returns false if user has no active role (sql.ErrNoRows)
//   - Returns false if blocked_until has expired (block is over)
//   - Returns false if userRole is nil
//
// Database Operations:
//   - SELECT from user_roles - via GetActiveUserRoleByUserID adapter
//   - Uses provided transaction
//
// Side Effects:
//   - None (read-only operation)
//   - Logs ERROR if operations fail
//
// Error Handling:
//   - sql.ErrNoRows is NOT an error (user has no active role, returns false)
//   - Infrastructure errors logged as ERROR
//   - Errors marked in span for distributed tracing
//   - Returns generic InternalError to prevent information disclosure
//
// Observability:
//   - Log entry on error: "permission.user.temp_block.check.db_failed" (ERROR)
//   - Span error marking for failed operations
//
// Important Notes:
//   - MUST be called within an existing transaction
//   - Caller MUST manage transaction lifecycle (commit/rollback)
//   - Preferred over IsUserTempBlocked when inside signin flow (already has tx)
//
// Example (Signin Flow with Transaction):
//
//	tx, _ := globalService.StartTransaction(ctx)
//	defer rollbackOnError(tx)
//
//	isBlocked, err := userService.IsUserTempBlockedWithTx(ctx, tx, userID)
//	if err != nil {
//	    return  // Infrastructure error
//	}
//	if isBlocked {
//	    return AuthenticationError("Invalid credentials")  // Generic error
//	}
//
//	// Continue signin flow...
func (us *userService) IsUserTempBlockedWithTx(ctx context.Context, tx *sql.Tx, userID int64) (bool, error) {
	// Initialize tracing for observability
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Retrieve user's active role
	userRole, err := us.repo.GetActiveUserRoleByUserID(ctx, tx, userID)
	if err != nil {
		// No active role is not an error (user not blocked)
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.Error("permission.user.temp_block.check.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}

	// Safety check: nil role means not blocked
	if userRole == nil {
		return false, nil
	}

	// Check if user is temp blocked and if the block hasn't expired
	if userRole.GetStatus() == globalmodel.StatusTempBlocked {
		blockedUntil := userRole.GetBlockedUntil()
		if blockedUntil != nil && time.Now().UTC().Before(*blockedUntil) {
			return true, nil
		}
	}

	return false, nil
}
