package userservices

import (
	"context"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user by setting their status to StatusTempBlocked
func (u *userService) BlockUserTemporarily(ctx context.Context, userID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	// Update all user roles to StatusTempBlocked
	err = u.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(permissionmodel.StatusTempBlocked))
	if err != nil {
		slog.ErrorContext(ctx, "failed to block user temporarily",
			"userID", userID,
			"error", err.Error())
		return utils.InternalError(fmt.Sprintf("failed to block user %d temporarily", userID))
	}

	// Clear user permissions cache to force refresh on next signin attempt
	err = u.permissionService.ClearUserPermissionsCache(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "failed to clear user permissions cache after blocking",
			"userID", userID,
			"error", err.Error())
		// Don't return error here as the blocking was successful
	}

	slog.InfoContext(ctx, "user blocked temporarily due to failed signin attempts",
		"userID", userID)

	return nil
}

// UnblockUserTemporarily unblocks a user by setting their status back to StatusActive
func (u *userService) UnblockUserTemporarily(ctx context.Context, userID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	// Update all user roles to StatusActive
	err = u.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(permissionmodel.StatusActive))
	if err != nil {
		slog.ErrorContext(ctx, "failed to unblock user temporarily",
			"userID", userID,
			"error", err.Error())
		return utils.InternalError(fmt.Sprintf("failed to unblock user %d", userID))
	}

	// Reset wrong signin attempts counter
	err = u.repo.ResetUserWrongSigninAttempts(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to reset wrong signin attempts after unblocking",
			"userID", userID,
			"error", err.Error())
		return utils.InternalError(fmt.Sprintf("failed to reset signin attempts for user %d", userID))
	}

	// Clear user permissions cache to force refresh
	err = u.permissionService.ClearUserPermissionsCache(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "failed to clear user permissions cache after unblocking",
			"userID", userID,
			"error", err.Error())
		// Don't return error here as the unblocking was successful
	}

	slog.InfoContext(ctx, "user unblocked temporarily",
		"userID", userID)

	return nil
}
