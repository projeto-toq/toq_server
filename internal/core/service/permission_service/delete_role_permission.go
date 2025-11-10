package permissionservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteRolePermission remove uma associação role-permission.
func (p *permissionServiceImpl) DeleteRolePermission(ctx context.Context, rolePermissionID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if rolePermissionID <= 0 {
		return utils.BadRequest("invalid role permission id")
	}

	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role_permission.delete.tx_start_failed", "role_permission_id", rolePermissionID, "error", txErr)
		return utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.role_permission.delete.tx_rollback_failed", "role_permission_id", rolePermissionID, "error", rbErr)
			}
		}
	}()

	existing, repoErr := p.permissionRepository.GetRolePermissionByID(ctx, tx, rolePermissionID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.role_permission.delete.repo_get_error", "role_permission_id", rolePermissionID, "error", repoErr)
		opErr = utils.InternalError("")
		return opErr
	}
	if existing == nil {
		opErr = utils.NotFoundError("role permission")
		return opErr
	}

	affectedUsers, usersErr := p.permissionRepository.GetActiveUserIDsByRoleID(ctx, tx, existing.GetRoleID())
	if usersErr != nil {
		utils.SetSpanError(ctx, usersErr)
		logger.Error("permission.role_permission.delete.get_active_users_failed", "role_id", existing.GetRoleID(), "error", usersErr)
		opErr = utils.InternalError("")
		return opErr
	}

	if delErr := p.permissionRepository.DeleteRolePermission(ctx, tx, rolePermissionID); delErr != nil {
		utils.SetSpanError(ctx, delErr)
		logger.Error("permission.role_permission.delete.repo_error", "role_permission_id", rolePermissionID, "error", delErr)
		opErr = utils.InternalError("")
		return opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.role_permission.delete.tx_commit_failed", "role_permission_id", rolePermissionID, "error", commitErr)
		return utils.InternalError("")
	}

	for _, uid := range affectedUsers {
		p.InvalidateUserCacheSafe(ctx, uid, "delete_role_permission")
	}

	logger.Info("permission.role_permission.deleted", "role_permission_id", rolePermissionID)
	return nil
}
