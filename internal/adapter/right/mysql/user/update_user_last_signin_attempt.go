package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserLastSignInAttempt updates only the last_signin_attempt timestamp for a specific user
//
// This function performs a targeted update of a single field, avoiding the overhead of updating
// all user fields. Used during authentication flow to record the timestamp of account lockout
// due to excessive failed signin attempts.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency with signin flow)
//   - userID: ID of the user whose last_signin_attempt will be updated
//   - attemptTime: Timestamp to record (typically time.Now().UTC())
//
// Returns:
//   - error: sql.ErrNoRows if user not found or deleted, database errors
//
// Business Rules:
//   - Only updates users WHERE deleted = 0 (active users)
//   - Timestamp stored with microsecond precision (TIMESTAMP(6))
//   - Used for security monitoring and analytics
//
// Database Schema:
//   - Table: users
//   - Column: last_signin_attempt (TIMESTAMP(6) NULL)
//   - Condition: deleted = 0
//
// Edge Cases:
//   - Returns sql.ErrNoRows if user doesn't exist or is deleted
//   - NULL value cleared by passing zero time (time.Time{})
//   - Previous value overwritten (no history kept in this table)
//
// Performance:
//   - Single-field UPDATE (faster than full UpdateUserByID)
//   - Uses primary key index (id)
//   - Minimal transaction log impact
//
// Important Notes:
//   - Called only when user is blocked due to failed attempts
//   - Separate from temp_wrong_signin.last_attempt_at (tracks every failure)
//   - users.last_signin_attempt tracks ONLY the lockout moment
//
// Example:
//
//	// Record lockout timestamp when blocking user
//	lockoutTime := time.Now().UTC()
//	err := adapter.UpdateUserLastSignInAttempt(ctx, tx, userID, lockoutTime)
//	if err == sql.ErrNoRows {
//	    // User not found or already deleted
//	} else if err != nil {
//	    // Infrastructure error
//	}
func (ua *UserAdapter) UpdateUserLastSignInAttempt(ctx context.Context, tx *sql.Tx, userID int64, attemptTime time.Time) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query updates only last_signin_attempt field
	// WHERE ensures user exists and is not soft-deleted
	query := `UPDATE users 
	          SET last_signin_attempt = ? 
	          WHERE id = ? AND deleted = 0`

	// Execute update via instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "update", query, attemptTime, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_last_signin_attempt.exec_error",
			"user_id", userID, "attempt_time", attemptTime, "error", execErr)
		return fmt.Errorf("update user last_signin_attempt: %w", execErr)
	}

	// Check if any rows were affected (does user exist and is active?)
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.update_last_signin_attempt.rows_affected_error",
			"user_id", userID, "error", raErr)
		return fmt.Errorf("get rows affected: %w", raErr)
	}

	// Return sql.ErrNoRows if user not found or deleted
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
