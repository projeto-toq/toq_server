package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RemoveDeviceTokensByDeviceID removes all push notification tokens for a specific device
//
// This function performs a bulk deletion of all device tokens associated with a specific device.
// Used during single-device logout (signout) or when a device session is revoked.
// Returns no error if device has no tokens (0 rows affected is considered success).
//
// Parameters:
//   - ctx: Context for tracing and logging propagation
//   - tx: Database transaction (can be nil for standalone operation)
//   - userID: User's unique identifier (scopes the deletion to prevent cross-user operations)
//   - deviceID: Unique device identifier (UUIDv4)
//
// Returns:
//   - error: Database errors (connection, query execution failures)
//
// Business Rules:
//   - Removes ALL tokens for the specific device (scoped by user_id + device_id)
//   - 0 rows affected is NOT an error (device may have no tokens)
//   - User scoping prevents accidental cross-user deletion
//
// Use Cases:
//   - Single-device logout: remove tokens only from current device
//   - Session revocation: cleanup when refresh token is revoked
//   - Device unregister: remove tokens when device is no longer used
//
// Example:
//
//	err := adapter.RemoveDeviceTokensByDeviceID(ctx, tx, 123, "550e8400-e29b-41d4-a716-446655440000")
//	// Returns nil even if device had no tokens
func (ua *UserAdapter) RemoveDeviceTokensByDeviceID(ctx context.Context, tx *sql.Tx, userID int64, deviceID string) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// DELETE query scoped by user_id AND device_id
	query := `DELETE FROM device_tokens WHERE user_id = ? AND device_id = ?`

	// Execute DELETE using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userID, deviceID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.remove_device_tokens_by_device.exec_error", "user_id", userID, "device_id", deviceID, "error", execErr)
		return fmt.Errorf("remove device tokens by device id: %w", execErr)
	}

	// Log number of tokens removed (informational, not error if 0)
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		// Non-critical: we can't log rows affected, but deletion succeeded
		logger.Warn("mysql.user.remove_device_tokens_by_device.rows_affected_warning", "user_id", userID, "device_id", deviceID, "error", raErr)
	} else if rowsAffected > 0 {
		logger.Debug("mysql.user.remove_device_tokens_by_device.success", "user_id", userID, "device_id", deviceID, "tokens_removed", rowsAffected)
	}

	// Success: 0 or more tokens removed
	return nil
}
