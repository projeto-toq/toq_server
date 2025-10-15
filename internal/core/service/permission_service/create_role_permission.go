package permissionservice

import (
	"context"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateRolePermission cria associação entre role e permissão.
func (p *permissionServiceImpl) CreateRolePermission(ctx context.Context, input CreateRolePermissionInput) (permissionmodel.RolePermissionInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.RoleID <= 0 {
		return nil, utils.ValidationError("roleId", "invalid role id")
	}
	if input.PermissionID <= 0 {
		return nil, utils.ValidationError("permissionId", "invalid permission id")
	}

	granted := true
	if input.Granted != nil {
		granted = *input.Granted
	}

	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role_permission.create.tx_start_failed", "role_id", input.RoleID, "permission_id", input.PermissionID, "error", txErr)
		return nil, utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.role_permission.create.tx_rollback_failed", "role_id", input.RoleID, "permission_id", input.PermissionID, "error", rbErr)
			}
		}
	}()

	if role, roleErr := p.permissionRepository.GetRoleByID(ctx, tx, input.RoleID); roleErr != nil || role == nil {
		if roleErr != nil {
			utils.SetSpanError(ctx, roleErr)
			logger.Error("permission.role_permission.create.get_role_failed", "role_id", input.RoleID, "error", roleErr)
			opErr = utils.InternalError("")
		} else {
			opErr = utils.NotFoundError("role")
		}
		return nil, opErr
	}

	permission, permErr := p.permissionRepository.GetPermissionByID(ctx, tx, input.PermissionID)
	if permErr != nil {
		utils.SetSpanError(ctx, permErr)
		logger.Error("permission.role_permission.create.get_permission_failed", "permission_id", input.PermissionID, "error", permErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}
	if permission == nil {
		opErr = utils.NotFoundError("permission")
		return nil, opErr
	}

	existing, existingErr := p.permissionRepository.GetRolePermissionByRoleIDAndPermissionID(ctx, tx, input.RoleID, input.PermissionID)
	if existingErr != nil {
		utils.SetSpanError(ctx, existingErr)
		logger.Error("permission.role_permission.create.get_existing_failed", "role_id", input.RoleID, "permission_id", input.PermissionID, "error", existingErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}
	if existing != nil {
		opErr = utils.ConflictError("role permission already exists")
		return nil, opErr
	}

	rolePermission := permissionmodel.NewRolePermission()
	rolePermission.SetRoleID(input.RoleID)
	rolePermission.SetPermissionID(input.PermissionID)
	rolePermission.SetGranted(granted)
	if input.Conditions != nil {
		rolePermission.SetConditions(input.Conditions)
	}

	affectedUsers, usersErr := p.permissionRepository.GetActiveUserIDsByRoleID(ctx, tx, input.RoleID)
	if usersErr != nil {
		utils.SetSpanError(ctx, usersErr)
		logger.Error("permission.role_permission.create.get_active_users_failed", "role_id", input.RoleID, "error", usersErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	if err = p.permissionRepository.CreateRolePermission(ctx, tx, rolePermission); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("permission.role_permission.create.repo_error", "role_id", input.RoleID, "permission_id", input.PermissionID, "error", err)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.role_permission.create.tx_commit_failed", "role_permission_id", rolePermission.GetID(), "error", commitErr)
		return nil, utils.InternalError("")
	}

	for _, uid := range affectedUsers {
		p.invalidateUserCacheSafe(ctx, uid, "create_role_permission")
	}

	logger.Info("permission.role_permission.created", "role_permission_id", rolePermission.GetID(), "role_id", input.RoleID, "permission_id", input.PermissionID)
	return rolePermission, nil
}
