package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteValidation removes all temporary validation codes for a specific user
//
// This function deletes the temp_user_validations record for a user, removing
// all pending validation codes (email, phone, password). This is typically
// called after successful validation or when explicitly revoking all pending
// validations for security reasons.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for standalone cleanup)
//   - id: User ID whose validation record will be deleted
//
// Returns:
//   - deleted: Number of rows deleted (0 or 1)
//   - error: Database errors (nil if record not found - idempotent operation)
//
// Business Rules:
//   - **Idempotent operation**: Returns success (nil error) even if no record exists
//   - Does NOT return sql.ErrNoRows (unlike other delete methods)
//   - Justification: Validation cleanup is considered "success" even if nothing to clean
//   - Service layer treats this as fire-and-forget operation
//
// Database Constraints:
//   - FK: user_id REFERENCES users(id) ON DELETE CASCADE
//
// Usage Example:
//
//	deleted, err := adapter.DeleteValidation(ctx, tx, userID)
//	// err == nil even if deleted == 0 (idempotent behavior)
func (ua *UserAdapter) DeleteValidation(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete validation record for the user
	query := `DELETE FROM temp_user_validations WHERE user_id = ?;`

	// Execute deletion using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_validation.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete validation: %w", execErr)
	}

	// Retrieve number of deleted rows
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_validation.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete validation rows affected: %w", rowsErr)
	}

	// Idempotent: if no rows were deleted, that's fine (nothing to clean up)
	// ✅ Does NOT return sql.ErrNoRows - this is intentional design
	return rowsAffected, nil
}
