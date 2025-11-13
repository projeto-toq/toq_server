package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RemoveDeviceToken deletes a specific push notification token by token string
//
// This function removes a single device token matching the provided token string.
// Returns sql.ErrNoRows if the token doesn't exist for the user.
//
// Parameters:
//   - ctx: Context for tracing and logging propagation
//   - tx: Database transaction (can be nil for standalone operation)
//   - userID: User's unique identifier (scopes the deletion)
//   - token: The FCM/APNs token string to remove
//
// Returns:
//   - error: sql.ErrNoRows if token not found, or database errors
//
// Business Rules:
//   - Only removes tokens belonging to the specified user (prevents cross-user deletion)
//   - Returns sql.ErrNoRows if token not found (allows service layer to map to 404)
//
// Example:
//
//	err := adapter.RemoveDeviceToken(ctx, tx, 123, "fcm_token_abc...")
//	if errors.Is(err, sql.ErrNoRows) {
//	    // Token not found (already removed or never existed)
//	}
func (ua *UserAdapter) RemoveDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// DELETE query scoped by user_id and device_token
	query := `DELETE FROM device_tokens WHERE user_id = ? AND device_token = ?`

	// Execute DELETE using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userID, token)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.remove_device_token.exec_error", "user_id", userID, "error", execErr)
		return fmt.Errorf("remove device token: %w", execErr)
	}

	// Check if any rows were affected (token existed and was deleted)
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.remove_device_token.rows_affected_error", "user_id", userID, "error", raErr)
		return fmt.Errorf("get rows affected: %w", raErr)
	}

	// Return sql.ErrNoRows if token not found (service layer maps to 404)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
