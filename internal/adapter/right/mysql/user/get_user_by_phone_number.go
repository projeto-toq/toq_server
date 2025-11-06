package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserByPhoneNumber retrieves a user by their phone number
//
// This function searches for a user by their unique phone number (E.164 format).
// Unlike GetUserByID, this method does NOT filter by deleted status, returning
// even soft-deleted users. This behavior is intentional for security checks.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - phoneNumber: User's phone in E.164 format (e.g., "+5511999999999")
//
// Returns:
//   - user: UserInterface with all user data
//   - error: sql.ErrNoRows if not found, or database errors
//
// Business Rules:
//   - Phone number has UNIQUE constraint in database
//   - Does NOT filter by deleted status (returns deleted users too)
//   - Service layer decides if deleted users are acceptable
//   - Phone format validation is service layer responsibility
//
// Query Logic:
//   - SELECT all user columns WHERE phone_number = ?
//   - No deleted filter (security/authentication checks need to know about deleted accounts)
//   - Exact string match (no normalization at DB level)
//
// Edge Cases:
//   - Phone not found: Returns sql.ErrNoRows
//   - User is soft-deleted (deleted=1): Still returned
//   - Multiple users with same phone: Prevented by UNIQUE constraint (DB error if occurs)
//   - Phone format variations: Exact match required (service normalizes to E.164)
//
// Performance:
//   - Uses UNIQUE index on phone_number column for fast lookup
//   - Single table query, no JOINs
//
// Security Considerations:
//   - Used for authentication flows (sign-in, password reset)
//   - Returns deleted users to prevent account enumeration attacks
//   - Service layer logs authentication attempts for deleted accounts
//   - PII data - ensure proper logging controls
//
// Deleted User Handling:
//   - Authentication: Service returns "Invalid credentials" (same as wrong password)
//   - Password reset: Service may allow reset to enable account recovery
//   - Sign-up: Service checks deleted status and may offer account reactivation
//
// Use Cases:
//   - User authentication (sign-in, refresh token)
//   - Phone number uniqueness validation during registration
//   - Password reset flow
//   - Phone number verification during profile update
//   - Security checks (detecting banned/deleted accounts attempting access)
//
// Example:
//
//	user, err := adapter.GetUserByPhoneNumber(ctx, tx, "+5511999999999")
//	if err == sql.ErrNoRows {
//	    // No user exists with this phone number
//	}
//	if user.IsDeleted() {
//	    // Account was soft-deleted, handle appropriately
//	}
func (ua *UserAdapter) GetUserByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (user usermodel.UserInterface, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query by phone_number (unique constraint) - includes deleted users for security
	// Note: No deleted filter allows authentication checks to detect deleted accounts
	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE phone_number = ?`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, phoneNumber)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_by_phone.query_error", "error", queryErr)
		return nil, fmt.Errorf("get user by phone number query: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entities
	entities, err := scanUserEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_phone.scan_error", "error", err)
		return nil, fmt.Errorf("scan user by phone rows: %w", err)
	}

	// Handle no results: phone number not registered
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: unique constraint should prevent multiple rows
	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple users found for phone number: %s", phoneNumber)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_phone.multiple_users_error", "phone_number", phoneNumber, "error", errMultiple)
		return nil, errMultiple
	}

	// Convert database entity to domain model
	user = userconverters.UserEntityToDomain(entities[0])

	return user, nil
}
