package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ExistsPhoneForAnotherUser checks if a phone is already used by a different user (deleted = 0)
//
// This function validates phone uniqueness during profile updates, ensuring users cannot
// change their phone to one already registered by another active user.
//
// Query Logic:
//   - Searches for active users (deleted = 0) with matching phone number
//   - Excludes specified user ID (allows user to keep their own phone)
//   - Returns boolean indicating if conflict exists
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - phone: Phone number to check for uniqueness (E.164 format expected)
//   - excludeUserID: User ID to exclude from check (typically current user)
//
// Returns:
//   - exists: true if phone is used by another active user
//   - error: Database connection errors or query execution failures
//
// Business Rules:
//   - Only checks active users (deleted = 0)
//   - Excludes specified user ID (user can keep their own phone unchanged)
//   - Exact string comparison (no normalization applied at DB level)
//   - Service layer should return 409 Conflict if exists=true
//
// Use Cases:
//   - User profile update: changing phone to new number
//   - Admin panel: validating phone before manual update
//   - Phone change confirmation: validating before sending SMS verification code
//
// Edge Cases:
//   - excludeUserID = 0: Checks all users (useful for registration validation)
//   - phone empty: Query executes but unlikely to match
//   - Multiple matches: Returns true (count > 0)
//   - Phone format variations: Service must normalize to E.164 before calling
//
// Performance:
//   - Uses UNIQUE index on phone_number column for fast lookup
//   - COUNT() aggregation avoids fetching full row data
//
// Example:
//
//	exists, err := adapter.ExistsPhoneForAnotherUser(ctx, tx,
//	    "+5511988887777",
//	    currentUserID)
//	if err != nil {
//	    // Handle infrastructure error
//	}
//	if exists {
//	    // Return 409 Conflict: phone already in use by another user
//	}
func (ua *UserAdapter) ExistsPhoneForAnotherUser(ctx context.Context, tx *sql.Tx, phone string, excludeUserID int64) (bool, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query counts active users with matching phone number, excluding specified user ID
	// Note: deleted = 0 ensures only active accounts are checked
	// Note: id <> ? allows user to keep their current phone unchanged
	query := `SELECT COUNT(id) as cnt FROM users WHERE phone_number = ? AND id <> ? AND deleted = 0;`
	row := ua.QueryRowContext(ctx, tx, "select", query, phone, excludeUserID)

	// Scan count result (aggregate always returns one row, even if 0)
	var cnt int64
	if scanErr := row.Scan(&cnt); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.exists_phone_for_another.scan_error", "error", scanErr)
		return false, fmt.Errorf("exists phone for another user scan: %w", scanErr)
	}

	// Convert count to boolean (any match = phone exists for another user)
	return cnt > 0, nil
}
