package userservices

import (
	"context"
	"fmt"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user by setting their status to StatusTempBlocked
//
// This method updates all user roles to StatusTempBlocked and clears permissions cache.
// Used for temporary account suspensions (e.g., policy violations, security issues).
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - userID: ID of the user to block temporarily
//
// Returns:
//   - error: Infrastructure errors (database, cache) mapped to InternalError (500)
//
// Business Rules:
//   - Updates ALL user roles to StatusTempBlocked (not just active role)
//   - Does NOT set blocked_until timestamp (use BlockUserDueToFailedAttempts for that)
//   - Cache invalidation ensures user cannot authenticate until cache refresh
//
// Database Operations:
//   - UPDATE user_roles (status column) - via UpdateUserRoleStatusByUserID adapter
//
// Side Effects:
//   - Modifies user_roles table (all roles set to StatusTempBlocked)
//   - Invalidates Redis cache (permissions cleared)
//   - Logs INFO when block succeeds
//   - Logs ERROR if operations fail
//
// Error Handling:
//   - Infrastructure errors logged as ERROR with context
//   - Errors marked in span for distributed tracing
//   - Returns generic InternalError to prevent information disclosure
//
// Observability:
//   - Log entry on error: "user.block_temp.update_role_status_error"
//   - Log entry on cache error: "user.block_temp.clear_cache_failed" (WARN)
//   - Log entry on success: "user.block_temp.success" (INFO)
//   - Span error marking for failed operations
//
// Important Notes:
//   - Does NOT use transactions (standalone operation)
//   - Cache clearing failure is non-fatal (logged as WARN)
//   - For failed signin blocking, use blockUserDueToFailedAttempts instead
//
// Example (Admin Block):
//
//	err := userService.BlockUserTemporarily(ctx, userID)
//	if err != nil {
//	    return  // Error already logged
//	}
func (us *userService) BlockUserTemporarily(ctx context.Context, userID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.block_temp.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update all user roles to StatusTempBlocked
	err = us.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(globalmodel.StatusTempBlocked))
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.block_temp.update_role_status_error", "user_id", userID, "error", err)
		return utils.InternalError(fmt.Sprintf("failed to block user %d temporarily", userID))
	}

	// Clear user permissions cache to force refresh on next signin attempt
	err = us.permissionService.ClearUserPermissionsCache(ctx, userID)
	if err != nil {
		logger.Warn("user.block_temp.clear_cache_failed", "user_id", userID, "error", err)
		// Don't return error here as the blocking was successful
	}

	logger.Info("user.block_temp.success", "user_id", userID)

	return nil
}
