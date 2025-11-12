package mysqldevicetokenadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListTokensByUserIDIfOptedIn retrieves FCM device tokens for a specific user if they opted in for notifications
//
// This function is used to send targeted push notifications to a single user.
// Only returns tokens for users who have opt_status = 1 (opted in) and are not deleted.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - userID: The unique identifier of the user
//
// Returns:
//   - tokens: Slice of distinct FCM device tokens (empty if user opted out or deleted)
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
//   - Only sends to users who explicitly opted in
//   - Respects soft-delete (no notifications to deleted accounts)
func (a *DeviceTokenAdapter) ListTokensByUserIDIfOptedIn(ctx context.Context, tx *sql.Tx, userID int64) ([]string, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query device tokens for user with opt-in and not deleted
	// INNER JOIN ensures only tokens for existing users are returned
	query := `SELECT DISTINCT dt.device_token 
			  FROM device_tokens dt 
			  INNER JOIN users u ON dt.user_id = u.id 
			  WHERE u.id = ? AND u.opt_status = 1 AND u.deleted = 0`

	// Execute query using instrumented adapter
	rows, queryErr := a.QueryContext(ctx, tx, "select", query, userID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.device_token.list_by_user.query_error", "user_id", userID, "error", queryErr)
		return nil, fmt.Errorf("list device tokens by user id: %w", queryErr)
	}
	defer rows.Close()

	tokens := make([]string, 0)
	for rows.Next() {
		var token string
		if scanErr := rows.Scan(&token); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.device_token.list_by_user.scan_error", "user_id", userID, "error", scanErr)
			return nil, fmt.Errorf("scan device token: %w", scanErr)
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.device_token.list_by_user.rows_error", "user_id", userID, "error", err)
		return nil, fmt.Errorf("iterate device tokens: %w", err)
	}

	return tokens, nil
}
