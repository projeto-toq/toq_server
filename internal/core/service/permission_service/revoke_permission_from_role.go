package permissionservice

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RevokePermissionFromRole revoga uma permissão de um role
func (p *permissionServiceImpl) RevokePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	if permissionID <= 0 {
		return utils.BadRequest("invalid permission id")
	}

	logger.Debug("permission.role_permission.revoke.request", "role_id", roleID, "permission_id", permissionID)

	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role_permission.revoke.tx_start_failed", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.role_permission.revoke.tx_rollback_failed", "role_id", roleID, "permission_id", permissionID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	rolePermission, err := p.permissionRepository.GetRolePermissionByRoleIDAndPermissionID(ctx, tx, roleID, permissionID)
	if err != nil {
		logger.Error("permission.role_permission.revoke.db_failed", "stage", "get_relation", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if rolePermission == nil {
		return utils.NotFoundError("role permission")
	}

	affectedUserIDs, err := p.permissionRepository.GetActiveUserIDsByRoleID(ctx, tx, roleID)
	if err != nil {
		logger.Error("permission.role_permission.revoke.get_role_users_failed", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	if err = p.permissionRepository.DeleteRolePermission(ctx, tx, rolePermission.GetID()); err != nil {
		logger.Error("permission.role_permission.revoke.db_failed", "stage", "delete_relation", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.role_permission.revoke.tx_commit_failed", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	logger.Info("permission.permission.revoked", "role_id", roleID, "permission_id", permissionID)
	for _, uid := range affectedUserIDs {
		p.invalidateUserCacheSafe(ctx, uid, "revoke_permission_from_role")
	}
	return nil
}
