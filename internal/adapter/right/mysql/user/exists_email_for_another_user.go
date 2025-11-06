package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ExistsEmailForAnotherUser checks if an email is already used by a different user (deleted = 0)
//
// This function validates email uniqueness during profile updates, ensuring users cannot
// change their email to one already registered by another active user.
//
// Query Logic:
//   - Searches for active users (deleted = 0) with matching email
//   - Excludes specified user ID (allows user to keep their own email)
//   - Returns boolean indicating if conflict exists
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - email: Email address to check for uniqueness
//   - excludeUserID: User ID to exclude from check (typically current user)
//
// Returns:
//   - exists: true if email is used by another active user
//   - error: Database connection errors or query execution failures
//
// Business Rules:
//   - Only checks active users (deleted = 0)
//   - Excludes specified user ID (user can keep their own email unchanged)
//   - Case-sensitive comparison (email collation determines behavior)
//   - Service layer should return 409 Conflict if exists=true
//
// Use Cases:
//   - User profile update: changing email to new address
//   - Admin panel: validating email before manual update
//   - Email change confirmation: validating before sending verification code
//
// Edge Cases:
//   - excludeUserID = 0: Checks all users (useful for registration validation)
//   - email empty: Query executes but unlikely to match
//   - Multiple matches: Returns true (count > 0)
//
// Performance:
//   - Uses UNIQUE index on email column for fast lookup
//   - COUNT() aggregation avoids fetching full row data
//
// Example:
//
//	exists, err := adapter.ExistsEmailForAnotherUser(ctx, tx,
//	    "newemail@example.com",
//	    currentUserID)
//	if err != nil {
//	    // Handle infrastructure error
//	}
//	if exists {
//	    // Return 409 Conflict: email already in use by another user
//	}
func (ua *UserAdapter) ExistsEmailForAnotherUser(ctx context.Context, tx *sql.Tx, email string, excludeUserID int64) (bool, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query counts active users with matching email, excluding specified user ID
	// Note: deleted = 0 ensures only active accounts are checked
	// Note: id <> ? allows user to keep their current email unchanged
	query := `SELECT COUNT(id) as cnt FROM users WHERE email = ? AND id <> ? AND deleted = 0;`
	row := ua.QueryRowContext(ctx, tx, "select", query, email, excludeUserID)

	// Scan count result (aggregate always returns one row, even if 0)
	var cnt int64
	if scanErr := row.Scan(&cnt); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.exists_email_for_another.scan_error", "error", scanErr)
		return false, fmt.Errorf("exists email for another user scan: %w", scanErr)
	}

	// Convert count to boolean (any match = email exists for another user)
	return cnt > 0, nil
}
