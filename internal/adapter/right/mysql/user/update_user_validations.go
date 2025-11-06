package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserValidations inserts or updates validation codes for email, phone, and password reset flows
//
// This function uses UPSERT (INSERT ... ON DUPLICATE KEY UPDATE) to either create a new
// validation record or update an existing one for the user. Used during:
//   - Email change flow (stores new email + verification code)
//   - Phone change flow (stores new phone + verification code)
//   - Password reset flow (stores password reset code)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency with user updates)
//   - validation: ValidationInterface with user ID and validation codes to store
//
// Returns:
//   - error: Database errors, constraint violations
//
// Business Rules:
//   - Primary key: user_id (max 1 validation record per user)
//   - UPSERT always affects exactly 1 row (INSERT or UPDATE)
//   - Does NOT delete when all codes are empty (caller should call DeleteValidation explicitly)
//   - Expiration timestamps must be set by service layer (typically NOW() + TTL)
//
// Database Schema:
//   - Table: temp_user_validations
//   - Columns: user_id (PK), new_email, email_code, email_code_exp,
//     new_phone, phone_code, phone_code_exp, password_code, password_code_exp
//   - All code fields are nullable (allows partial updates)
//
// Edge Cases:
//   - Empty validation codes: Stored as NULL (service layer should validate before calling)
//   - Expired codes: NOT cleaned up by this function (use worker or DeleteExpiredValidations)
//   - User doesn't exist: Foreign key constraint prevents insert (returns error)
//
// Performance:
//   - Single-row UPSERT using PRIMARY KEY (very fast)
//   - ON DUPLICATE KEY UPDATE avoids SELECT + INSERT/UPDATE race condition
//
// Important:
//   - Never returns sql.ErrNoRows (UPSERT always affects 1 row)
//   - To delete validation record, call DeleteValidation() explicitly
//   - Service layer responsible for clearing codes when no longer needed
//
// Example:
//
//	validation := usermodel.NewValidation()
//	validation.SetUserID(userID)
//	validation.SetEmailCode("123456")
//	validation.SetEmailCodeExp(time.Now().Add(15 * time.Minute))
//
//	err := adapter.UpdateUserValidations(ctx, tx, validation)
//	if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) UpdateUserValidations(ctx context.Context, tx *sql.Tx, validation usermodel.ValidationInterface) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// UPSERT query: insert new validation record or update existing one
	// Note: PRIMARY KEY on user_id ensures max 1 record per user
	// Note: ON DUPLICATE KEY UPDATE replaces all columns (not partial)
	query := `INSERT INTO temp_user_validations (
		user_id, new_email, email_code, email_code_exp, 
		new_phone, phone_code, phone_code_exp, password_code, password_code_exp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		new_email = VALUES(new_email),
		email_code = VALUES(email_code),
		email_code_exp = VALUES(email_code_exp),
		new_phone = VALUES(new_phone),
		phone_code = VALUES(phone_code),
		phone_code_exp = VALUES(phone_code_exp),
		password_code = VALUES(password_code),
		password_code_exp = VALUES(password_code_exp)`

	// Convert domain model to database entity
	entity := userconverters.UserValidationDomainToEntity(validation)

	// Execute UPSERT using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "insert", query,
		entity.UserID,
		entity.NewEmail,
		entity.EmailCode,
		entity.EmailCodeExp,
		entity.NewPhone,
		entity.PhoneCode,
		entity.PhoneCodeExp,
		entity.PasswordCode,
		entity.PasswordCodeExp,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_validations.exec_error", "user_id", entity.UserID, "error", execErr)
		return fmt.Errorf("update user validations: %w", execErr)
	}

	// Verify operation success (UPSERT should always affect at least 1 row)
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.update_user_validations.rows_affected_error", "user_id", entity.UserID, "error", raErr)
		return fmt.Errorf("user validations update rows affected: %w", raErr)
	}

	// Safety check: UPSERT should always affect exactly 1 row
	if rowsAffected == 0 {
		errZeroRows := fmt.Errorf("upsert affected 0 rows (unexpected)")
		utils.SetSpanError(ctx, errZeroRows)
		logger.Error("mysql.user.update_user_validations.zero_rows_error", "user_id", entity.UserID, "error", errZeroRows)
		return errZeroRows
	}

	return nil
}
