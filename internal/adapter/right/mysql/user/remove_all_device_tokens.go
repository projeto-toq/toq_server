package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RemoveAllDeviceTokensByUserID removes all push notification tokens for a user
//
// This function performs a bulk deletion of all device tokens associated with a user.
// Used during logout (global signout), account deletion, or opt-out from notifications.
// Returns no error if user has no tokens (0 rows affected is considered success).
//
// Parameters:
//   - ctx: Context for tracing and logging propagation
//   - tx: Database transaction (can be nil for standalone operation)
//   - userID: User's unique identifier
//
// Returns:
//   - error: Database errors (connection, query execution failures)
//
// Business Rules:
//   - Removes ALL tokens for the user (no filtering by device_id)
//   - 0 rows affected is NOT an error (user may have no tokens)
//   - Should be called within transaction when part of larger operation (e.g., account deletion)
//
// Use Cases:
//   - Global logout: remove all tokens from all devices
//   - Account deletion: cleanup before removing user record
//   - Opt-out: remove all tokens when user disables notifications
//
// Example:
//
//	err := adapter.RemoveAllDeviceTokensByUserID(ctx, tx, 123)
//	// Returns nil even if user had no tokens
func (ua *UserAdapter) RemoveAllDeviceTokensByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// DELETE query scoped by user_id (removes all tokens)
	query := `DELETE FROM device_tokens WHERE user_id = ?`

	// Execute DELETE using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.remove_all_device_tokens.exec_error", "user_id", userID, "error", execErr)
		return fmt.Errorf("remove all device tokens: %w", execErr)
	}

	// Log number of tokens removed (informational, not error if 0)
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		// Non-critical: we can't log rows affected, but deletion succeeded
		logger.Warn("mysql.user.remove_all_device_tokens.rows_affected_warning", "user_id", userID, "error", raErr)
	} else if rowsAffected > 0 {
		logger.Debug("mysql.user.remove_all_device_tokens.success", "user_id", userID, "tokens_removed", rowsAffected)
	}

	// Success: 0 or more tokens removed
	return nil
}
