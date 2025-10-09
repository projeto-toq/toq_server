package userservices

import (
	"context"
	"fmt"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user by setting their status to StatusTempBlocked
func (u *userService) BlockUserTemporarily(ctx context.Context, userID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.block_temp.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update all user roles to StatusTempBlocked
	err = u.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(permissionmodel.StatusTempBlocked))
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.block_temp.update_role_status_error", "user_id", userID, "error", err)
		return utils.InternalError(fmt.Sprintf("failed to block user %d temporarily", userID))
	}

	// Clear user permissions cache to force refresh on next signin attempt
	err = u.permissionService.ClearUserPermissionsCache(ctx, userID)
	if err != nil {
		logger.Warn("user.block_temp.clear_cache_failed", "user_id", userID, "error", err)
		// Don't return error here as the blocking was successful
	}

	logger.Info("user.block_temp.success", "user_id", userID)

	return nil
}

// UnblockUserTemporarily unblocks a user by setting their status back to StatusActive
func (u *userService) UnblockUserTemporarily(ctx context.Context, userID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.unblock_temp.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update all user roles to StatusActive
	err = u.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(permissionmodel.StatusActive))
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.unblock_temp.update_role_status_error", "user_id", userID, "error", err)
		return utils.InternalError(fmt.Sprintf("failed to unblock user %d", userID))
	}

	// Reset wrong signin attempts counter
	err = u.repo.ResetUserWrongSigninAttempts(ctx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.unblock_temp.reset_wrong_signin_error", "user_id", userID, "error", err)
		return utils.InternalError(fmt.Sprintf("failed to reset signin attempts for user %d", userID))
	}

	// Clear user permissions cache to force refresh
	err = u.permissionService.ClearUserPermissionsCache(ctx, userID)
	if err != nil {
		logger.Warn("user.unblock_temp.clear_cache_failed", "user_id", userID, "error", err)
		// Don't return error here as the unblocking was successful
	}

	logger.Info("user.unblock_temp.success", "user_id", userID)

	return nil
}
