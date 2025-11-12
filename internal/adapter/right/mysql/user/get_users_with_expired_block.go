package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
)

// GetUsersWithExpiredBlock returns users whose temporary block has expired
//
// This function retrieves all users with blocked_until <= NOW().
// Used by worker process to automatically unblock users whose block period has ended.
//
// NEW ARCHITECTURE:
//   - Queries users table directly (not user_roles)
//   - Returns UserInterface (not UserRoleInterface)
//   - Simpler query (no JOIN required)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for read-only queries)
//
// Returns:
//   - users: Slice of UserInterface with expired blocks (empty slice if none)
//   - error: Database errors, scan errors
//
// Business Rules:
//   - Returns only non-deleted users (WHERE deleted = 0)
//   - Filters by blocked_until IS NOT NULL AND blocked_until <= NOW()
//   - Does NOT modify data (read-only query)
//
// Query Structure:
//   - Table: users
//   - Filter: deleted = 0 AND blocked_until IS NOT NULL AND blocked_until <= NOW()
//   - Columns: all user columns (uses rowsToUserEntities helper)
//
// Edge Cases:
//   - No expired blocks: Returns empty slice (NOT sql.ErrNoRows)
//   - blocked_until = NULL: Excluded from results
//   - Multiple expired blocks: All returned (should be processed by worker)
//
// Performance:
//   - Uses idx_users_blocking index (blocked_until, permanently_blocked)
//   - Called periodically by worker (every 5 minutes)
//
// Example:
//
//	// Worker periodic check
//	expiredUsers, err := adapter.GetUsersWithExpiredBlock(ctx, nil)
//	if err != nil {
//	    logger.Error("Failed to get expired blocks", "error", err)
//	    return
//	}
//	for _, user := range expiredUsers {
//	    // Unblock each user
//	    _ = service.ClearUserBlockedUntil(ctx, tx, user.GetID())
//	}
func (ua *UserAdapter) GetUsersWithExpiredBlock(ctx context.Context, tx *sql.Tx) ([]usermodel.UserInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query users whose block period has expired
	// Note: blocked_until IS NOT NULL ensures only valid temp blocks are returned
	// Note: blocked_until <= NOW() identifies expired blocks
	query := `
		SELECT id, full_name, nick_name, national_id, creci_number, creci_state, 
		       creci_validity, born_at, phone_number, email, zip_code, street, 
		       number, complement, neighborhood, city, state, password, 
		       opt_status, last_activity_at, deleted, blocked_until, permanently_blocked
		FROM users
		WHERE deleted = 0 
		  AND blocked_until IS NOT NULL 
		  AND blocked_until <= NOW()
	`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_users_with_expired_block.query_error", "error", queryErr)
		return nil, fmt.Errorf("query users with expired block: %w", queryErr)
	}
	defer rows.Close()

	var users []usermodel.UserInterface

	// Scan each row into UserEntity
	for rows.Next() {
		var entity userentity.UserEntity

		if err := rows.Scan(
			&entity.ID,
			&entity.FullName,
			&entity.NickName,
			&entity.NationalID,
			&entity.CreciNumber,
			&entity.CreciState,
			&entity.CreciValidity,
			&entity.BornAt,
			&entity.PhoneNumber,
			&entity.Email,
			&entity.ZipCode,
			&entity.Street,
			&entity.Number,
			&entity.Complement,
			&entity.Neighborhood,
			&entity.City,
			&entity.State,
			&entity.Password,
			&entity.OptStatus,
			&entity.LastActivityAt,
			&entity.Deleted,
			&entity.BlockedUntil,
			&entity.PermanentlyBlocked,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.user.get_users_with_expired_block.scan_error", "error", err)
			return nil, fmt.Errorf("scan user with expired block: %w", err)
		}

		// Convert entity to domain model
		user := userconverters.UserEntityToDomain(entity)
		users = append(users, user)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_users_with_expired_block.rows_error", "error", err)
		return nil, fmt.Errorf("iterate users with expired block: %w", err)
	}

	logger.Debug("mysql.user.get_users_with_expired_block.success", "count", len(users))
	return users, nil
}
