package mysqluseradapter

import (
	"context"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ResetUserWrongSigninAttempts deletes failed signin attempt tracking for a specific user
//
// This function removes the wrong signin tracking record from temp_wrong_signin table,
// effectively resetting the failed attempt counter to zero. Used after successful
// authentication or when admin manually unlocks an account.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - userID: User ID whose failed signin attempts should be reset
//
// Returns:
//   - error: Database errors (does NOT return error if record doesn't exist)
//
// Business Rules:
//   - Deletes tracking record by user_id (PRIMARY KEY)
//   - No error if record doesn't exist (0 rows affected is success)
//   - Used after successful signin to clear previous failed attempts
//   - Used by admin to manually unlock temporarily blocked accounts
//
// Database Schema:
//   - Table: temp_wrong_signin
//   - Primary Key: user_id
//   - Columns: user_id, failed_attempts, last_attempt_at
//
// Edge Cases:
//   - User never had failed attempts: DELETE affects 0 rows (not an error)
//   - Record already deleted: DELETE affects 0 rows (not an error)
//   - Invalid user ID: DELETE affects 0 rows (not an error)
//
// Performance:
//   - Single-row DELETE using PRIMARY KEY (very fast)
//   - No WHERE clause needed beyond user_id
//
// Important Notes:
//   - Does NOT validate if user exists in users table
//   - Does NOT check user's deleted status
//   - Transaction parameter intentionally nil (standalone operation for performance)
//
// Example:
//
//	// After successful signin
//	err := adapter.ResetUserWrongSigninAttempts(ctx, userID)
//	if err != nil {
//	    // Log error but don't fail signin (non-critical cleanup)
//	    logger.Warn("Failed to reset wrong signin attempts", "error", err)
//	}
func (ua *UserAdapter) ResetUserWrongSigninAttempts(ctx context.Context, userID int64) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
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
	result, execErr := ua.ExecContext(ctx, nil, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.reset_wrong_signin_attempts.exec_error", "user_id", userID, "error", execErr)
		return fmt.Errorf("reset wrong signin attempts: %w", execErr)
	}

	// Check rows affected (not treated as error if 0, record may not exist)
	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.reset_wrong_signin_attempts.rows_affected_error", "user_id", userID, "error", rowsErr)
		return fmt.Errorf("reset wrong signin attempts rows affected: %w", rowsErr)
	}

	return nil
}
