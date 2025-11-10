package permissionservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeletePermission remove uma permissão de forma definitiva.
func (p *permissionServiceImpl) DeletePermission(ctx context.Context, permissionID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if permissionID <= 0 {
		return utils.BadRequest("invalid permission id")
	}

	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.permission.delete.tx_start_failed", "permission_id", permissionID, "error", txErr)
		return utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.permission.delete.tx_rollback_failed", "permission_id", permissionID, "error", rbErr)
			}
		}
	}()

	existing, repoErr := p.permissionRepository.GetPermissionByID(ctx, tx, permissionID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.permission.delete.repo_get_error", "permission_id", permissionID, "error", repoErr)
		opErr = utils.InternalError("")
		return opErr
	}
	if existing == nil {
		opErr = utils.NotFoundError("permission")
		return opErr
	}

	roleIDs, roleErr := p.permissionRepository.GetRoleIDsByPermissionID(ctx, tx, permissionID)
	if roleErr != nil {
		utils.SetSpanError(ctx, roleErr)
		logger.Error("permission.permission.delete.get_role_ids_failed", "permission_id", permissionID, "error", roleErr)
		opErr = utils.InternalError("")
		return opErr
	}

	userIDSet := make(map[int64]struct{})
	for _, roleID := range roleIDs {
		ids, usersErr := p.permissionRepository.GetActiveUserIDsByRoleID(ctx, tx, roleID)
		if usersErr != nil {
			utils.SetSpanError(ctx, usersErr)
			logger.Error("permission.permission.delete.get_role_users_failed", "permission_id", permissionID, "role_id", roleID, "error", usersErr)
			opErr = utils.InternalError("")
			return opErr
		}
		for _, uid := range ids {
			userIDSet[uid] = struct{}{}
		}
	}

	if delErr := p.permissionRepository.DeletePermission(ctx, tx, permissionID); delErr != nil {
		utils.SetSpanError(ctx, delErr)
		logger.Error("permission.permission.delete.repo_error", "permission_id", permissionID, "error", delErr)
		opErr = utils.InternalError("")
		return opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.permission.delete.tx_commit_failed", "permission_id", permissionID, "error", commitErr)
		return utils.InternalError("")
	}

	for uid := range userIDSet {
		p.InvalidateUserCacheSafe(ctx, uid, "delete_permission")
	}

	logger.Info("permission.deleted", "permission_id", permissionID)
	return nil
}
