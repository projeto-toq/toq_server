package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserLastActivity updates the last_activity_at timestamp for a single user
//
// This function records the current UTC timestamp as the user's last activity time.
// Used for idle timeout calculations, analytics, and user session management.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for standalone updates)
//   - id: User ID whose last activity timestamp should be updated
//
// Returns:
//   - error: Database errors (does NOT return error if user not found)
//
// Business Rules:
//   - Always sets timestamp to current UTC time (NOW())
//   - Updates ANY user regardless of deleted status
//   - No validation of user existence (0 rows affected is not an error)
//
// Performance:
//   - Single-row update using PRIMARY KEY (very fast)
//   - Called frequently (on every authenticated request)
//   - Consider using batch update (BatchUpdateUserLastActivity) for high-traffic scenarios
//
// Use Cases:
//   - Record user activity after each authenticated API call
//   - Calculate idle timeout for automatic logout
//   - Analytics: track user engagement patterns
//   - Admin dashboard: show "last seen" timestamp
//
// Important Notes:
//   - For high-traffic apps, use BatchUpdateUserLastActivity instead
//   - Timestamp stored in UTC (convert to user timezone in presentation layer)
//   - Does NOT check deleted status (updates even deleted users)
//
// Example:
//
//	err := adapter.UpdateUserLastActivity(ctx, nil, userID)
//	if err != nil {
//	    // Log error but don't fail request (non-critical operation)
//	    logger.Warn("Failed to update last activity", "error", err)
//	}
func (ua *UserAdapter) UpdateUserLastActivity(ctx context.Context, tx *sql.Tx, id int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update last_activity_at to current UTC timestamp
	// Note: No WHERE deleted = 0 check - updates ANY user (active or deleted)
	query := `UPDATE users SET last_activity_at = ? WHERE id = ?;`

	// Execute update using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		time.Now().UTC(),
		id,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_last_activity.exec_error", "user_id", id, "error", execErr)
		return fmt.Errorf("update user last_activity: %w", execErr)
	}

	// Check rows affected (not treated as error if 0, user may not exist)
	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_last_activity.rows_affected_error", "user_id", id, "error", rowsErr)
		return fmt.Errorf("user last activity update rows affected: %w", rowsErr)
	}

	return
}
