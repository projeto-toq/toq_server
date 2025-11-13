package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListDeviceTokenStringsByOptedInUsers retrieves FCM device tokens for all users who opted in for notifications
//
// This function is used for broadcast/bulk notifications sent to all opted-in users.
// Only returns tokens for users with opt_status = 1 (opted in) and not deleted.
// Returns empty slice if no users opted in (NOT sql.ErrNoRows).
//
// Parameters:
//   - ctx: Context for tracing and logging propagation
//   - tx: Database transaction (can be nil for standalone queries)
//
// Returns:
//   - tokens: Slice of distinct FCM token strings from all opted-in users
//   - error: Database query errors
//
// Business Rules:
//   - Filters by opt_status = 1 (user consented to notifications)
//   - Filters by deleted = 0 (excludes soft-deleted users)
//   - Returns DISTINCT tokens (prevents duplicate sends)
//   - Returns empty slice if no users opted in (NOT sql.ErrNoRows)
//
// Performance Considerations:
//   - Can return large result set (thousands of tokens)
//   - Consider pagination for very large user bases
//   - Uses INNER JOIN to avoid orphaned tokens
//
// Use Cases:
//   - System-wide announcements (new features, maintenance)
//   - Marketing campaigns (respects opt-in consent)
//   - Emergency alerts
//
// Security:
//   - Only sends to users who explicitly opted in (GDPR compliant)
//   - Respects soft-delete (no notifications to deleted accounts)
//
// Example:
//
//	tokens, err := adapter.ListDeviceTokenStringsByOptedInUsers(ctx, nil)
//	// len(tokens) == 0 if no users opted in (not an error)
func (ua *UserAdapter) ListDeviceTokenStringsByOptedInUsers(ctx context.Context, tx *sql.Tx) ([]string, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query all device tokens for users with opt-in and not deleted
	// INNER JOIN ensures only tokens for existing users are returned
	query := `SELECT DISTINCT dt.device_token 
	          FROM device_tokens dt 
	          INNER JOIN users u ON dt.user_id = u.id 
	          WHERE u.opt_status = 1 AND u.deleted = 0`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.list_device_tokens_all_opted_in.query_error", "error", queryErr)
		return nil, fmt.Errorf("list device tokens for opted-in users: %w", queryErr)
	}
	defer rows.Close()

	// Scan token strings
	tokens := make([]string, 0)
	for rows.Next() {
		var token string
		if scanErr := rows.Scan(&token); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.user.list_device_tokens_all_opted_in.scan_error", "error", scanErr)
			return nil, fmt.Errorf("scan device token: %w", scanErr)
		}
		tokens = append(tokens, token)
	}

	// Check for iteration errors
	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.list_device_tokens_all_opted_in.rows_error", "error", rowsErr)
		return nil, fmt.Errorf("iterate device tokens: %w", rowsErr)
	}

	return tokens, nil
}
