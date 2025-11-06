package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteWrongSignInByUserID removes failed sign-in attempt tracking for a user
//
// This function deletes the temp_wrong_signin record for a user, clearing
// the counter of failed authentication attempts. This is typically called
// after a successful sign-in or when manually resetting security lockouts.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - security operations must be transactional)
//   - userID: ID of the user whose failed attempts will be cleared
//
// Returns:
//   - deleted: Number of rows deleted (0 or 1)
//   - error: sql.ErrNoRows if user has no failed attempts record, database errors otherwise
//
// Business Rules:
//   - Returns sql.ErrNoRows if user has no tracking record (never failed sign-in)
//   - Service layer maps sql.ErrNoRows to success (no failed attempts = success)
//   - Resets lockout counter and allows user to attempt sign-in again
//
// Database Constraints:
//   - FK: user_id REFERENCES users(id) ON DELETE CASCADE
//   - Unique: user_id (one record per user maximum)
//
// Usage Example:
//
//	deleted, err := adapter.DeleteWrongSignInByUserID(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    // User never had failed attempts (expected scenario)
//	} else if err != nil {
//	    // Infrastructure error
//	}
func (ua *UserAdapter) DeleteWrongSignInByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete failed sign-in tracking record
	query := `DELETE FROM temp_wrong_signin WHERE user_id = ?;`

	// Execute deletion using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_wrong_signin.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete temp_wrong_signin: %w", execErr)
	}

	// Check if tracking record was found and deleted
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_wrong_signin.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete temp_wrong_signin rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user had no failed attempts record
	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return rowsAffected, nil
}
