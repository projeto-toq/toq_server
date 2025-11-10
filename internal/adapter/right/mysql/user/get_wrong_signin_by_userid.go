package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetWrongSigninByUserID retrieves failed sign-in attempt tracking for a specific user
//
// This function fetches the temporary record from temp_wrong_signin table that tracks
// failed authentication attempts for rate limiting and account protection. Each user
// has at most one tracking record (user_id is PRIMARY KEY).
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - id: User's unique identifier
//
// Returns:
//   - wrongSignin: WrongSigninInterface with failed attempts count and timestamp
//   - error: sql.ErrNoRows if no tracking record exists, or database errors
//
// Business Rules:
//   - user_id is PRIMARY KEY in temp_wrong_signin (max 1 record per user)
//   - Record created on first failed authentication attempt
//   - Record deleted after successful sign-in or when timeout expires
//   - Service layer implements rate limiting logic based on failed_attempts
//
// Table Schema (temp_wrong_signin):
//   - user_id: INT (PRIMARY KEY, FOREIGN KEY to users.id)
//   - failed_attempts: INT NOT NULL - Count of consecutive failed sign-ins
//   - last_attempt_at: TIMESTAMP NOT NULL - Timestamp of most recent failed attempt
//
// Edge Cases:
//   - User has no tracking record: Returns sql.ErrNoRows (normal state, no failed attempts)
//   - Multiple records for same user: Returns error (PRIMARY KEY should prevent this)
//   - Record with failed_attempts = 0: Valid state (just reset, awaiting cleanup)
//
// Performance:
//   - PRIMARY KEY lookup (user_id) - extremely fast
//   - Single row maximum (no pagination needed)
//
// Security Considerations:
//   - Used for brute force protection
//   - Threshold typically 3-5 failed attempts before temporary block
//   - Timeout period typically 15-30 minutes
//   - Service layer enforces lockout logic
//
// Use Cases:
//   - Sign-in flow: Check if user is rate-limited before validating password
//   - Account security: Display failed attempt count to user
//   - Admin panel: Monitor suspicious authentication patterns
//
// Example:
//
//	wrongSignin, err := adapter.GetWrongSigninByUserID(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    // No failed attempts on record - user can sign in
//	}
//	if wrongSignin.GetFailedAttempts() >= 5 {
//	    // User is rate-limited due to excessive failed attempts
//	}
func (ua *UserAdapter) GetWrongSigninByUserID(ctx context.Context, tx *sql.Tx, id int64) (usermodel.WrongSigninInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query wrong signin tracking by user ID
	// Note: Primary key is user_id, so max 1 row expected
	query := `SELECT user_id, failed_attempts, last_attempt_at 
	          FROM temp_wrong_signin WHERE user_id = ?`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, id)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_wrong_signin.query_error", "user_id", id, "error", queryErr)
		return nil, fmt.Errorf("get wrong signin by user id query: %w", queryErr)
	}
	defer rows.Close()

	// Scan rows using type-safe function
	entities, err := scanWrongSigninEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_wrong_signin.scan_error", "user_id", id, "error", err)
		return nil, fmt.Errorf("scan wrong signin rows: %w", err)
	}

	// Handle no results: user has no failed attempts on record
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: primary key should prevent multiple rows
	if len(entities) > 1 {
		errMultiple := errors.New("multiple wrong_signin rows found for user")
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_wrong_signin.multiple_rows_error", "user_id", id, "count", len(entities), "error", errMultiple)
		return nil, errMultiple
	}

	// Convert entity to domain model using type-safe converter
	wrongSignin := userconverters.WrongSignInEntityToDomainTyped(entities[0])

	return wrongSignin, nil
}
