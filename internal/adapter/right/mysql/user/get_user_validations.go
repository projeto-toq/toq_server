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

// GetUserValidations retrieves temporary validation data for a user
//
// This function fetches pending validation records from temp_user_validations table,
// which stores temporary codes and new values for email/phone changes and password resets.
// Each user has at most one validation record (user_id is PRIMARY KEY).
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - id: User's unique identifier
//
// Returns:
//   - validation: ValidationInterface with codes, expiration times, and new values
//   - error: sql.ErrNoRows if no validation record exists, or database errors
//
// Business Rules:
//   - user_id is PRIMARY KEY in temp_user_validations (max 1 record per user)
//   - All code fields are nullable (only populated for active validation flows)
//   - Expiration times are nullable (NULL if no active validation for that channel)
//   - Record exists only during validation flows (deleted after completion)
//
// Table Schema (temp_user_validations):
//   - user_id: INT (PRIMARY KEY, FOREIGN KEY to users.id)
//   - new_email: VARCHAR(100) NULL - New email pending verification
//   - email_code: VARCHAR(10) NULL - Verification code for email change
//   - email_code_exp: TIMESTAMP NULL - Expiration time for email code
//   - new_phone: VARCHAR(25) NULL - New phone pending verification
//   - phone_code: VARCHAR(10) NULL - Verification code for phone change
//   - phone_code_exp: TIMESTAMP NULL - Expiration time for phone code
//   - password_code: VARCHAR(10) NULL - Verification code for password reset
//   - password_code_exp: TIMESTAMP NULL - Expiration time for password code
//
// Nullable Field Handling:
//   - NULL values converted to empty strings/zero times in domain model
//   - Service layer checks for empty values to determine active validations
//   - Expired codes are NOT automatically deleted (cleanup job handles this)
//
// Edge Cases:
//   - User has no validation record: Returns sql.ErrNoRows (normal state)
//   - All fields NULL except user_id: Returns valid validation with empty codes
//   - Expired codes still in DB: Returned as-is, service validates expiration
//   - Multiple validation types active: Single record holds all (email + phone + password)
//
// Performance:
//   - PRIMARY KEY lookup (user_id) - extremely fast
//   - Single row maximum (no pagination needed)
//
// Validation Flow Integration:
//  1. User initiates change (email/phone) or password reset
//  2. Service creates/updates record in temp_user_validations
//  3. Verification code sent to user
//  4. User submits code
//  5. Service calls GetUserValidations to verify code and check expiration
//  6. On success: applies change and deletes validation record
//  7. On failure: validation record remains for retry (until expiration or max attempts)
//
// Security Considerations:
//   - Codes are time-limited (typically 5-15 minutes)
//   - PII data (new_email, new_phone) - ensure proper logging controls
//   - Rate limiting should be applied at service layer
//
// Use Cases:
//   - Email change verification: Check email_code against user input
//   - Phone change verification: Check phone_code against SMS input
//   - Password reset: Verify password_code before allowing new password
//   - Multi-factor authentication: Validate temporary codes
//
// Example:
//
//	validation, err := adapter.GetUserValidations(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    // No pending validations for this user
//	}
//	if validation.GetEmailCode() != "" {
//	    // User has pending email change validation
//	    if time.Now().After(validation.GetEmailCodeExp()) {
//	        // Code expired
//	    }
//	}
func (ua *UserAdapter) GetUserValidations(ctx context.Context, tx *sql.Tx, id int64) (validation usermodel.ValidationInterface, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query validation record by user ID
	// Note: user_id is PRIMARY KEY, so maximum 1 row expected
	query := `SELECT user_id, new_email, email_code, email_code_exp, 
	          new_phone, phone_code, phone_code_exp, 
	          password_code, password_code_exp 
	          FROM temp_user_validations WHERE user_id = ?;`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, id)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_validations.query_error", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	// Scan rows using type-safe function
	entities, err := scanValidationEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_validations.scan_error", "error", err)
		return nil, fmt.Errorf("scan user validations rows: %w", err)
	}

	// Handle no results: user has no pending validation flows
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: PRIMARY KEY should prevent multiple rows
	if len(entities) > 1 {
		errMultiple := errors.New("multiple validations found for user")
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_validations.multiple_validations_error", "user_id", id, "error", errMultiple)
		return nil, errMultiple
	}

	// Convert entity to domain model using type-safe converter
	validation = userconverters.UserValidationEntityToDomainTyped(entities[0])

	return
}
