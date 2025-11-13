package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListDeviceTokenStringsByUserIDIfOptedIn retrieves FCM device token strings for a specific user if they opted in for notifications
//
// This function is used to send targeted push notifications to a single user.
// Only returns tokens for users who have opt_status = 1 (opted in) and are not deleted.
// Returns empty slice if user opted out or is deleted (NOT sql.ErrNoRows).
//
// Parameters:
//   - ctx: Context for tracing and logging propagation
//   - tx: Database transaction (can be nil for standalone queries)
//   - userID: The unique identifier of the user
//
// Returns:
//   - tokens: Slice of distinct FCM token strings (empty if user opted out or deleted)
//   - error: Database query errors
//
// Business Rules:
//   - Filters by opt_status = 1 (user consented to notifications)
//   - Filters by deleted = 0 (excludes soft-deleted users)
//   - Returns DISTINCT tokens (user may have multiple devices with same token)
//   - Returns empty slice if user opted out or is deleted (NOT sql.ErrNoRows)
//
// Use Cases:
//   - Sending personalized notifications (e.g., listing updates, messages)
//   - Triggering alerts for specific user actions
//
// Security:
//   - Only sends to users who explicitly opted in (GDPR compliant)
//   - Respects soft-delete (no notifications to deleted accounts)
//
// Example:
//
//	tokens, err := adapter.ListDeviceTokenStringsByUserIDIfOptedIn(ctx, nil, 123)
//	// len(tokens) == 0 if user opted out (not an error)
func (ua *UserAdapter) ListDeviceTokenStringsByUserIDIfOptedIn(ctx context.Context, tx *sql.Tx, userID int64) ([]string, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query device tokens for user with opt-in and not deleted
	// INNER JOIN ensures only tokens for existing users are returned
	query := `SELECT DISTINCT dt.device_token 
	          FROM device_tokens dt 
	          INNER JOIN users u ON dt.user_id = u.id 
	          WHERE u.id = ? AND u.opt_status = 1 AND u.deleted = 0`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, userID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.list_device_tokens_opted_in.query_error", "user_id", userID, "error", queryErr)
		return nil, fmt.Errorf("list device tokens by user id if opted in: %w", queryErr)
	}
	defer rows.Close()

	// Scan token strings
	tokens := make([]string, 0)
	for rows.Next() {
		var token string
		if scanErr := rows.Scan(&token); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.user.list_device_tokens_opted_in.scan_error", "user_id", userID, "error", scanErr)
			return nil, fmt.Errorf("scan device token: %w", scanErr)
		}
		tokens = append(tokens, token)
	}

	// Check for iteration errors
	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.list_device_tokens_opted_in.rows_error", "user_id", userID, "error", rowsErr)
		return nil, fmt.Errorf("iterate device tokens: %w", rowsErr)
	}

	return tokens, nil
}
