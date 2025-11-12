package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ResetUserWrongSigninAttempts deletes failed signin attempt tracking for a specific user
//
// This function removes the wrong signin tracking record from temp_wrong_signin table,
// effectively resetting the failed attempt counter to zero. Used after successful
// authentication or when admin/worker manually unlocks an account.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for atomic unblock operations)
//   - userID: User ID whose failed signin attempts should be reset
//
// Returns:
//   - error: Database errors (does NOT return error if record doesn't exist)
//
// Business Rules:
//   - Deletes tracking record by user_id (PRIMARY KEY)
//   - No error if record doesn't exist (0 rows affected is success)
//   - Used after successful signin to clear previous failed attempts
//   - Used by worker/admin to unlock temporarily blocked accounts
//   - MUST run within transaction for consistency with unblock operations
//
// Database Schema:
//   - Table: temp_wrong_signin
//   - Primary Key: user_id
//   - ON DELETE CASCADE: automatically removed when user deleted
//
// Edge Cases:
//   - Record doesn't exist: Success (0 rows affected, no error)
//   - Multiple calls: Idempotent (second call succeeds with 0 rows affected)
//
// Performance:
//   - PRIMARY KEY deletion (extremely fast)
//   - Single row maximum
//
// Important Notes:
//   - Changed from standalone method to transactional (breaking change)
//   - Previously accepted no tx parameter, now REQUIRES tx for consistency
//   - Caller MUST manage transaction lifecycle (commit/rollback)
//
// Example:
//
//	// Clear failed attempts after successful signin
//	tx, _ := globalService.StartTransaction(ctx)
//	defer rollbackOnError(tx)
//
//	err := adapter.ResetUserWrongSigninAttempts(ctx, tx, userID)
//
//	_ = globalService.CommitTransaction(ctx, tx)
func (ua *UserAdapter) ResetUserWrongSigninAttempts(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete wrong signin tracking record by user ID
	// Note: PRIMARY KEY ensures max 1 row deleted
	// Note: 0 rows affected is not an error (record may not exist)
	query := `DELETE FROM temp_wrong_signin WHERE user_id = ?`

	// Execute DELETE using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.reset_wrong_signin_attempts.exec_error",
			"user_id", userID, "error", execErr)
		return fmt.Errorf("reset user wrong signin attempts: %w", execErr)
	}

	// Check rows affected (informational only, not an error if 0)
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.reset_wrong_signin_attempts.rows_affected_error",
			"user_id", userID, "error", raErr)
		return fmt.Errorf("get rows affected: %w", raErr)
	}

	// Success regardless of rows affected (idempotent operation)
	if rowsAffected > 0 {
		logger.Debug("mysql.user.reset_wrong_signin_attempts.success",
			"user_id", userID, "rows_deleted", rowsAffected)
	} else {
		logger.Debug("mysql.user.reset_wrong_signin_attempts.no_record",
			"user_id", userID)
	}

	return nil
}
