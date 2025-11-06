package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteExpiredValidations removes validation records where all codes are expired or empty
//
// This function performs cleanup of temp_user_validations table by deleting rows
// where ALL THREE validation codes (email, phone, password) are either:
//   - NULL
//   - Empty string (”)
//   - Expired (expiration timestamp < NOW())
//
// The deletion is limited to prevent long-running transactions and table locks.
// This function should be called periodically by a cleanup worker.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for standalone cleanup)
//   - limit: Maximum number of rows to delete in one execution (recommended: 1000)
//
// Returns:
//   - deleted: Number of rows actually deleted
//   - error: Database errors
//
// Business Rules:
//   - A row is deleted ONLY if ALL three code categories are invalid
//   - If ANY code is valid (not empty and not expired), row is preserved
//   - Uses AND logic to ensure no active validation is lost
//
// Performance:
//   - LIMIT clause prevents table lock on large datasets
//   - Should be called repeatedly until deleted = 0 for full cleanup
//   - Recommended interval: hourly via scheduled worker
//
// Example Usage:
//
//	deleted, err := adapter.DeleteExpiredValidations(ctx, nil, 1000)
//	for deleted > 0 && err == nil {
//	    deleted, err = adapter.DeleteExpiredValidations(ctx, nil, 1000)
//	}
func (ua *UserAdapter) DeleteExpiredValidations(ctx context.Context, tx *sql.Tx, limit int) (int64, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete only rows where ALL codes are expired/empty (AND logic)
	// This ensures we never delete a row with at least one active validation code
	//
	// Deletion criteria:
	//   email_code: (IS NULL OR = '' OR email_code_exp < NOW())
	//   AND phone_code: (IS NULL OR = '' OR phone_code_exp < NOW())
	//   AND password_code: (IS NULL OR = '' OR password_code_exp < NOW())
	//
	// LIMIT clause prevents long-running DELETE on large datasets
	query := `DELETE FROM temp_user_validations
		WHERE
			( (email_code IS NULL OR email_code = '' OR email_code_exp < NOW())
			AND (phone_code IS NULL OR phone_code = '' OR phone_code_exp < NOW())
			AND (password_code IS NULL OR password_code = '' OR password_code_exp < NOW()) )
		LIMIT ?;`

	// Execute deletion using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_expired_validations.exec_error", "limit", limit, "error", execErr)
		return 0, fmt.Errorf("delete expired validations: %w", execErr)
	}

	// Retrieve number of deleted rows for caller feedback
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_expired_validations.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete expired validations rows affected: %w", rowsErr)
	}

	// Log cleanup metrics (info level for scheduled operations)
	if rowsAffected > 0 {
		logger.Info("user.validations.cleanup.completed", "deleted_rows", rowsAffected, "limit", limit)
	}

	return rowsAffected, nil
}
