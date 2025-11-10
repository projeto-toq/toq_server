package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetExpiredTempBlockedUsers returns users whose temporary block has expired
//
// This function retrieves all active user roles with StatusTempBlocked where blocked_until
// timestamp has passed. Used by worker process to automatically restore access for users
// whose temporary block period has ended.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for read-only queries)
//
// Returns:
//   - userRoles: Slice of UserRoleInterface with expired blocks (empty slice if none)
//   - error: Database errors, scan errors
//
// Business Rules:
//   - Returns only active roles (WHERE ur.is_active = 1)
//   - Filters by StatusTempBlocked status
//   - Checks blocked_until <= NOW() (expired blocks)
//   - Does NOT modify data (read-only query)
//
// Query Structure:
//   - Table: user_roles
//   - Filter: status = StatusTempBlocked AND blocked_until IS NOT NULL AND blocked_until <= NOW()
//   - Columns: id, user_id, role_id, is_active, status, expires_at, blocked_until
//
// Edge Cases:
//   - No expired blocks: Returns empty slice (NOT sql.ErrNoRows)
//   - blocked_until = NULL: Excluded from results (malformed data)
//   - Multiple expired blocks: All returned (should be processed by worker)
//
// Performance:
//   - Table scan if no index on (status, blocked_until)
//   - Composite index on (status, blocked_until, is_active) recommended
//   - Called periodically by worker (every 1-5 minutes)
//
// Important Notes:
//   - Empty result is valid (no expired blocks to process)
//   - Worker calls UnblockUser() for each returned user role
//   - Should be called within worker's periodic task loop
//
// Example:
//
//	// Worker periodic check
//	expiredBlocks, err := adapter.GetExpiredTempBlockedUsers(ctx, nil)
//	if err != nil {
//	    logger.Error("Failed to get expired blocks", "error", err)
//	    return
//	}
//	for _, userRole := range expiredBlocks {
//	    // Unblock each user
//	    _ = service.UnblockUser(ctx, userRole.GetUserID())
//	}
func (ua *UserAdapter) GetExpiredTempBlockedUsers(ctx context.Context, tx *sql.Tx) ([]usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query active temp blocked users whose block period has expired
	// Note: blocked_until IS NOT NULL ensures only valid temp blocks are returned
	// Note: blocked_until <= NOW() identifies expired blocks
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.is_active, ur.status, ur.expires_at, ur.blocked_until
		FROM user_roles ur
		WHERE ur.status = ? 
		  AND ur.blocked_until IS NOT NULL 
		  AND ur.blocked_until <= NOW()
		  AND ur.is_active = 1
	`

	// Execute query using instrumented adapter (auto-generates metrics + tracing)
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, globalmodel.StatusTempBlocked)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_expired_temp_blocked_users.query_error", "error", queryErr)
		return nil, fmt.Errorf("query expired temp blocked users: %w", queryErr)
	}
	defer rows.Close()

	var userRoles []usermodel.UserRoleInterface

	// Scan each row into UserRoleInterface
	index := 0
	for rows.Next() {
		index++
		userRole := usermodel.NewUserRole()
		var (
			id           int64
			userID       int64
			roleID       int64
			isActive     bool
			status       globalmodel.UserRoleStatus
			expiresAt    sql.NullTime
			blockedUntil sql.NullTime
		)

		if err := rows.Scan(&id, &userID, &roleID, &isActive, &status, &expiresAt, &blockedUntil); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.user.get_expired_temp_blocked_users.scan_error", "row_index", index-1, "error", err)
			return nil, fmt.Errorf("scan expired temp blocked user row: %w", err)
		}

		// Populate user role domain model
		userRole.SetID(id)
		userRole.SetUserID(userID)
		userRole.SetRoleID(roleID)
		userRole.SetIsActive(isActive)
		userRole.SetStatus(status)

		if expiresAt.Valid {
			t := expiresAt.Time
			userRole.SetExpiresAt(&t)
		}
		if blockedUntil.Valid {
			t := blockedUntil.Time
			userRole.SetBlockedUntil(&t)
		}

		userRoles = append(userRoles, userRole)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_expired_temp_blocked_users.rows_error", "error", err)
		return nil, fmt.Errorf("iterate expired temp blocked users rows: %w", err)
	}

	logger.Debug("mysql.user.get_expired_temp_blocked_users.success", "count", len(userRoles))
	return userRoles, nil
}
