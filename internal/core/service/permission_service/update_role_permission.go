package permissionservice

import (
	"context"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateRolePermission atualiza campos permitidos de uma associação role-permission.
func (p *permissionServiceImpl) UpdateRolePermission(ctx context.Context, input UpdateRolePermissionInput) (permissionmodel.RolePermissionInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ID <= 0 {
		return nil, utils.BadRequest("invalid role permission id")
	}

	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role_permission.update.tx_start_failed", "role_permission_id", input.ID, "error", txErr)
		return nil, utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.role_permission.update.tx_rollback_failed", "role_permission_id", input.ID, "error", rbErr)
			}
		}
	}()

	existing, repoErr := p.permissionRepository.GetRolePermissionByID(ctx, tx, input.ID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.role_permission.update.repo_get_error", "role_permission_id", input.ID, "error", repoErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}
	if existing == nil {
		opErr = utils.NotFoundError("role permission")
		return nil, opErr
	}

	if input.Granted != nil {
		existing.SetGranted(*input.Granted)
	}

	affectedUsers, usersErr := p.permissionRepository.GetActiveUserIDsByRoleID(ctx, tx, existing.GetRoleID())
	if usersErr != nil {
		utils.SetSpanError(ctx, usersErr)
		logger.Error("permission.role_permission.update.get_active_users_failed", "role_id", existing.GetRoleID(), "error", usersErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	if err = p.permissionRepository.UpdateRolePermission(ctx, tx, existing); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("permission.role_permission.update.repo_error", "role_permission_id", input.ID, "error", err)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.role_permission.update.tx_commit_failed", "role_permission_id", input.ID, "error", commitErr)
		return nil, utils.InternalError("")
	}

	for _, uid := range affectedUsers {
		p.invalidateUserCacheSafe(ctx, uid, "update_role_permission")
	}

	logger.Info("permission.role_permission.updated", "role_permission_id", existing.GetID(), "granted", existing.GetGranted())
	return existing, nil
}
