package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUsersWithExpiredBlock fetches users with expired temporal blocks
//
// This method queries users table for rows where blocked_until <= NOW(),
// returning full UserInterface objects for processing by worker.
//
// Parameters:
//   - ctx: Context with logger and tracing
//   - tx: Database transaction (REQUIRED for consistent read)
//
// Returns:
//   - []UserInterface: Users with expired blocks (empty slice if none)
//   - error: Repository errors (database issues)
//
// Business Rules:
//   - Queries users.blocked_until <= NOW() AND deleted = 0
//   - Returns empty slice if no expired blocks found
//   - Used by TempBlockCleanerWorker every 5 minutes
//   - Worker calls ClearUserBlockedUntil for each returned user
//
// Database Operations:
//   - SELECT * FROM users WHERE blocked_until <= NOW() AND deleted = 0
//
// Side Effects:
//   - Read-only operation (no side effects)
//   - Logs ERROR if query fails
//
// Error Handling:
//   - Empty result is success (returns empty slice, no error)
//   - Infrastructure errors logged and returned
//
// Example:
//
//	tx, _ := globalService.StartTransaction(ctx)
//	users, err := us.GetUsersWithExpiredBlock(ctx, tx)
//	if err != nil {
//	    logger.Error("Failed to fetch expired blocks", "error", err)
//	    return
//	}
//	for _, user := range users {
//	    us.ClearUserBlockedUntil(ctx, tx, user.GetID())
//	}
func (us *userService) GetUsersWithExpiredBlock(ctx context.Context, tx *sql.Tx) ([]usermodel.UserInterface, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	users, err := us.repo.GetUsersWithExpiredBlock(ctx, tx)
	if err != nil {
		logger.Error("user_service.get_expired_blocks_failed", "error", err)
		return nil, err
	}

	logger.Debug("user_service.get_expired_blocks_success", "count", len(users))
	return users, nil
}
