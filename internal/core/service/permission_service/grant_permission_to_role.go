package permissionservice

import (
	"context"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GrantPermissionToRole concede uma permissão a um role
func (p *permissionServiceImpl) GrantPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	if permissionID <= 0 {
		return utils.BadRequest("invalid permission id")
	}

	slog.Debug("permission.permission.grant.start", "role_id", roleID, "permission_id", permissionID)

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("permission.permission.tx_start_failed", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("permission.permission.tx_rollback_failed", "role_id", roleID, "permission_id", permissionID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	// Verificar se o role existe
	role, err := p.permissionRepository.GetRoleByID(ctx, tx, roleID)
	if err != nil {
		slog.Error("permission.permission.get_role_failed", "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if role == nil {
		return utils.NotFoundError("Role")
	}

	// Verificar se a permissão existe
	permission, err := p.permissionRepository.GetPermissionByID(ctx, tx, permissionID)
	if err != nil {
		slog.Error("permission.permission.get_permission_failed", "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if permission == nil {
		return utils.NotFoundError("Permission")
	}

	// Verificar se a relação já existe
	existingRolePermission, err := p.permissionRepository.GetRolePermissionByRoleIDAndPermissionID(ctx, tx, roleID, permissionID)
	if err != nil {
		slog.Error("permission.permission.get_role_permission_failed", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if existingRolePermission != nil {
		return utils.ConflictError("role already has permission")
	}

	// Criar a nova RolePermission
	rolePermission := permissionmodel.NewRolePermission()
	rolePermission.SetRoleID(roleID)
	rolePermission.SetPermissionID(permissionID)
	rolePermission.SetGranted(true)

	// Salvar no banco
	err = p.permissionRepository.CreateRolePermission(ctx, tx, rolePermission)
	if err != nil {
		slog.Error("permission.permission.create_failed", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("permission.permission.tx_commit_failed", "role_id", roleID, "permission_id", permissionID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	slog.Info("permission.permission.granted", "role_id", roleID, "permission_id", permissionID)
	return nil
}
