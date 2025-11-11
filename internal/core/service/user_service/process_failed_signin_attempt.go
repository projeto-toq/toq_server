package userservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// processFailedSigninAttempt handles the logic for tracking and blocking users after failed signin attempts
//
// This method orchestrates the complete failed attempt handling flow:
//  1. Retrieves or initializes failed attempt tracking record
//  2. Increments failure counter and updates timestamp
//  3. Persists tracking record to database
//  4. Checks if threshold reached (MaxWrongSigninAttempts)
//  5. If threshold reached: blocks user temporarily and records lockout timestamp
//
// The function is called ONLY when password validation fails (bcrypt.CompareHashAndPassword != nil).
// It implements brute-force protection by temporarily locking accounts after excessive failures.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging (must contain logger)
//   - tx: Database transaction (REQUIRED for atomic tracking + blocking operations)
//   - userID: ID of the user who failed authentication
//
// Returns:
//   - error: Infrastructure errors (database, transaction) mapped to InternalError (500)
//
// Business Rules:
//   - Counter starts at 1 on first failure (record created via UPSERT)
//   - Each failure increments counter and updates last_attempt_at timestamp
//   - When counter reaches MaxWrongSigninAttempts (3), user is blocked temporarily
//   - Block duration: TempBlockDuration (15 minutes) from moment of lockout
//   - Lockout timestamp recorded in users.last_signin_attempt field
//   - Tracking record (temp_wrong_signin) persists until successful signin or timeout
//
// Database Operations:
//   - UPSERT temp_wrong_signin (tracking table)
//   - UPDATE user_roles (sets status to StatusTempBlocked) - via BlockUserTemporarily
//   - UPDATE users (records last_signin_attempt) - via UpdateUserLastSignInAttempt
//
// Side Effects:
//   - Modifies temp_wrong_signin table (counter incremented)
//   - May modify user_roles table (status change to blocked)
//   - May modify users table (last_signin_attempt timestamp set)
//   - Logs WARN when user is blocked (security event)
//   - Logs INFO on each failed attempt (observability)
//
// Error Handling:
//   - Infrastructure errors logged as ERROR and marked in span
//   - Returns InternalError to prevent information disclosure (no hints to attacker)
//   - Transaction rollback handled by caller (signIn function)
//
// Observability:
//   - Log entry on each failed attempt: "auth.signin.failed_attempt"
//   - Log entry when blocking: "auth.signin.user_blocked"
//   - Log entries for infrastructure errors
//   - Span error marking for distributed tracing
//
// Example Call Flow:
//
//	// In signIn function, after password validation fails:
//	if bcrypt.CompareHashAndPassword(...) != nil {
//	    err = us.processFailedSigninAttempt(ctx, tx, userID)
//	    if err != nil {
//	        return  // Error already logged and mapped
//	    }
//	    err = utils.AuthenticationError("Invalid credentials")
//	    return
//	}
func (us *userService) processFailedSigninAttempt(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Reuse context from parent (already has logger and tracer)
	// Do not start new tracer in private methods
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Retrieve existing tracking record or create new one
	wrongSignin, err := us.repo.GetWrongSigninByUserID(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// First failed attempt: initialize new tracking record
			wrongSignin = usermodel.NewWrongSignin()
		} else {
			// Infrastructure error retrieving tracking record
			logger.Error("auth.signin.wrong_signin_get_failed", "user_id", userID, "error", err)
			return utils.InternalError("Failed to check signin attempts")
		}
	}

	// Update tracking record: increment counter and set timestamp
	currentAttempts := wrongSignin.GetFailedAttempts()
	wrongSignin.SetUserID(userID)
	wrongSignin.SetLastAttemptAt(time.Now().UTC())
	wrongSignin.SetFailedAttempts(currentAttempts + 1)

	// Persist tracking record (UPSERT: INSERT if new, UPDATE if exists)
	err = us.repo.UpdateWrongSignIn(ctx, tx, wrongSignin)
	if err != nil {
		logger.Error("auth.signin.wrong_signin_update_failed", "user_id", userID, "error", err)
		return utils.InternalError("Failed to update signin attempts")
	}

	// Log the failed attempt for observability (not ERROR, it's expected behavior)
	logger.Info("auth.signin.failed_attempt",
		"security", true,
		"user_id", userID,
		"attempts", wrongSignin.GetFailedAttempts(),
		"max_attempts", usermodel.MaxWrongSigninAttempts)

	// Check if threshold reached: block user if at or above limit
	if wrongSignin.GetFailedAttempts() >= usermodel.MaxWrongSigninAttempts {
		// Block user temporarily (sets user_roles.status and blocked_until)
		err = us.blockUserDueToFailedAttempts(ctx, tx, userID)
		if err != nil {
			// Error already logged by blockUserDueToFailedAttempts
			return err
		}

		// Log security event: user blocked (WARN level for security monitoring)
		logger.Warn("auth.signin.user_blocked",
			"security", true,
			"user_id", userID,
			"attempts", wrongSignin.GetFailedAttempts(),
			"blocked_duration_minutes", int(usermodel.TempBlockDuration.Minutes()))
	}

	return nil
}
