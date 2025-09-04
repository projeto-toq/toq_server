package permissionservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RevokePermissionFromRole revoga uma permissão de um role
func (p *permissionServiceImpl) RevokePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	if permissionID <= 0 {
		return utils.BadRequest("invalid permission id")
	}

	slog.Debug("permission.role_permission.revoke.request", "role_id", roleID, "permission_id", permissionID)

	// Buscar a RolePermission existente
	rolePermission, err := p.permissionRepository.GetRolePermissionByRoleIDAndPermissionID(ctx, nil, roleID, permissionID)
	if err != nil {
		slog.Error("permission.role_permission.revoke.db_failed", "stage", "get_relation", "role_id", roleID, "permission_id", permissionID, "error", err)
		return utils.InternalError("")
	}
	if rolePermission == nil {
		return utils.NotFoundError("role permission")
	}

	// Remover a RolePermission
	err = p.permissionRepository.DeleteRolePermission(ctx, nil, rolePermission.GetID())
	if err != nil {
		slog.Error("permission.role_permission.revoke.db_failed", "stage", "delete_relation", "role_id", roleID, "permission_id", permissionID, "error", err)
		return utils.InternalError("")
	}

	slog.Info("permission.permission.revoked", "role_id", roleID, "permission_id", permissionID)
	return nil
}
