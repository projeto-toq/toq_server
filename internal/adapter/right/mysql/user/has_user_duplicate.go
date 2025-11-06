package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HasUserDuplicate checks if any active user exists with matching phone, email, or national ID
//
// This function prevents duplicate account creation by searching for existing users
// with at least one matching unique identifier. Used during registration and profile updates.
//
// Query Logic:
//   - Uses OR condition: match on ANY of phone/email/nationalID triggers duplicate
//   - Filters by deleted = 0 (only checks active users, ignoring soft-deleted accounts)
//   - Returns count > 0 as boolean (more efficient than SELECT *)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - user: UserInterface with phone, email, and national ID to check
//
// Returns:
//   - exists: true if ANY active user found with matching identifiers
//   - error: Database connection errors or query execution failures
//
// Business Rules:
//   - Soft-deleted users (deleted=1) are IGNORED and can be "reused"
//   - Service layer must handle exists=true by returning 409 Conflict error
//   - Used in: CreateUser, UpdateUserByID (when changing unique fields)
//
// Edge Cases:
//   - All three fields match same user: still returns true (count=1)
//   - Multiple users match different fields: returns true (count>1)
//   - Empty phone/email/nationalID: query still executes but unlikely to match
//
// Performance:
//   - Indexed columns (phone_number, email, national_id) ensure fast lookup
//   - COUNT() aggregation avoids fetching full row data
//
// Example:
//
//	user := usermodel.NewUser()
//	user.SetPhoneNumber("+5511999999999")
//	user.SetEmail("existing@example.com")
//	user.SetNationalID("12345678901")
//
//	exists, err := adapter.HasUserDuplicate(ctx, tx, user)
//	if err != nil {
//	    // Handle infrastructure error
//	}
//	if exists {
//	    // Return 409 Conflict: phone, email, or CPF already registered
//	}
func (ua *UserAdapter) HasUserDuplicate(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error) {
	// Initialize tracing for observability (distributed tracing + metrics)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query checks for any active user with matching phone, email, or national ID
	// Note: OR condition - match on ANY identifier triggers duplicate
	// Note: deleted = 0 filter excludes soft-deleted users (they can be "recycled")
	query := `SELECT count(id) as count
				FROM users WHERE (phone_number = ? OR email = ? OR national_id = ? ) AND deleted = 0`

	row := ua.QueryRowContext(ctx, tx, "select", query,
		user.GetPhoneNumber(),
		user.GetEmail(),
		user.GetNationalID(),
	)

	// Scan count result (aggregate always returns one row, even if 0)
	var qty int64
	if scanErr := row.Scan(&qty); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.has_duplicate.scan_error", "error", scanErr)
		return false, fmt.Errorf("check user duplicate: %w", scanErr)
	}

	// Convert count to boolean (any match = duplicate)
	exist = qty > 0

	return exist, nil
}
