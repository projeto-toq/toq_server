package userservices

import (
	"context"
	"errors"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetActiveRoleStatus returns only the status of the active user role for the current authenticated user.
// Domain invariant: an active role must always exist; its absence is treated as an infrastructure inconsistency.
func (us *userService) GetActiveRoleStatus(ctx context.Context) (status permissionmodel.UserRoleStatus, err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return 0, utils.BadRequest("Invalid user context")
	}

	tx, txErr := us.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("user.get_active_role_status.tx_start_error", "error", txErr)
		return 0, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("user.get_active_role_status.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	activeRole, aerr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if aerr != nil {
		utils.SetSpanError(ctx, aerr)
		slog.Error("user.get_active_role_status.read_active_role_error", "error", aerr, "user_id", userID)
		return 0, utils.InternalError("Failed to get active role")
	}
	if activeRole == nil {
		inconsistencyErr := errors.New("active role missing for user")
		utils.SetSpanError(ctx, inconsistencyErr)
		slog.Error("user.get_active_role_status.active_role_missing", "user_id", userID)
		return 0, utils.InternalError("Active role missing")
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		slog.Error("user.get_active_role_status.tx_commit_error", "error", cmErr)
		return 0, utils.InternalError("Failed to commit transaction")
	}

	return activeRole.GetStatus(), nil
}
