package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateWrongSignIn inserts or updates failed signin attempt tracking for a user
//
// This function uses UPSERT (INSERT ... ON DUPLICATE KEY UPDATE) to either create
// a new tracking record or increment the failed attempt counter. Used during
// authentication flow to track and limit consecutive failed signin attempts.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency with signin flow)
//   - wrongSignin: WrongSigninInterface with user ID, attempt count, and timestamp
//
// Returns:
//   - error: Database errors, constraint violations
//
// Business Rules:
//   - Primary key: user_id (max 1 tracking record per user)
//   - UPSERT always affects exactly 1 row (INSERT or UPDATE)
//   - Failed attempts counter incremented by service layer before calling
//   - Timestamp records last failed attempt time for rate limiting
//
// Database Schema:
//   - Table: temp_wrong_signin
//   - Primary Key: user_id (FK to users.id)
//   - Columns: user_id, failed_attempts, last_attempt_at
//   - ON DELETE CASCADE: Record deleted when user deleted
//
// Edge Cases:
//   - First failed attempt: INSERT creates new record with counter = 1
//   - Subsequent failures: UPDATE increments counter and updates timestamp
//   - User doesn't exist: Foreign key constraint prevents insert (returns error)
//
// Performance:
//   - Single-row UPSERT using PRIMARY KEY (very fast)
//   - ON DUPLICATE KEY UPDATE avoids SELECT + INSERT/UPDATE race condition
//   - Called on every failed signin attempt
//
// Important Notes:
//   - Never returns sql.ErrNoRows (UPSERT always affects 1 row)
//   - Reset via ResetUserWrongSigninAttempts() or DeleteWrongSignInByUserID()
//   - Service layer enforces max attempts limit and temporary blocking
//
// Example:
//
//	wrongSignin := usermodel.NewWrongSignin()
//	wrongSignin.SetUserID(userID)
//	wrongSignin.SetFailedAttempts(3) // Service incremented counter
//	wrongSignin.SetLastAttemptAt(time.Now())
//
//	err := adapter.UpdateWrongSignIn(ctx, tx, wrongSignin)
//	if err != nil {
//	    // Handle infrastructure error (don't block signin on tracking failure)
//	}
func (ua *UserAdapter) UpdateWrongSignIn(ctx context.Context, tx *sql.Tx, wrongSigin usermodel.WrongSigninInterface) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// UPSERT query: insert new tracking record or update existing one
	// Note: PRIMARY KEY on user_id ensures max 1 record per user
	// Note: ON DUPLICATE KEY UPDATE replaces counter and timestamp
	query := `INSERT INTO temp_wrong_signin (
				user_id, failed_attempts, last_attempt_at
				) VALUES (?, ?, ?)
				ON DUPLICATE KEY UPDATE
				failed_attempts = VALUES(failed_attempts),
				last_attempt_at = VALUES(last_attempt_at)`

	// Convert domain model to database entity
	entity := userconverters.WrongSignInDomainToEntity(wrongSigin)

	// Execute UPSERT using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "insert", query,
		entity.UserID,
		entity.FailedAttempts,
		entity.LastAttemptAT,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_wrong_signin.exec_error", "user_id", entity.UserID, "error", execErr)
		return fmt.Errorf("update wrong_signin: %w", execErr)
	}

	// Verify operation success (UPSERT should always affect 1 row)
	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_wrong_signin.rows_affected_error", "user_id", entity.UserID, "error", rowsErr)
		return fmt.Errorf("wrong_signin update rows affected: %w", rowsErr)
	}

	return nil
}
