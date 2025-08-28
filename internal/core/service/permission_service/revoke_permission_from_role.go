package permissionservice

import (
	"context"
	"fmt"
	"log/slog"
)

// RevokePermissionFromRole revoga uma permissão de um role
func (p *permissionServiceImpl) RevokePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	if roleID <= 0 {
		return fmt.Errorf("invalid role ID: %d", roleID)
	}

	if permissionID <= 0 {
		return fmt.Errorf("invalid permission ID: %d", permissionID)
	}

	slog.Debug("Revoking permission from role", "roleID", roleID, "permissionID", permissionID)

	// Buscar a RolePermission existente
	rolePermission, err := p.permissionRepository.GetRolePermissionByRoleIDAndPermissionID(ctx, nil, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to find role permission: %w", err)
	}
	if rolePermission == nil {
		return fmt.Errorf("role %d does not have permission %d", roleID, permissionID)
	}

	// Remover a RolePermission
	err = p.permissionRepository.DeleteRolePermission(ctx, nil, rolePermission.GetID())
	if err != nil {
		return fmt.Errorf("failed to revoke permission from role: %w", err)
	}

	slog.Info("Permission revoked from role successfully", "roleID", roleID, "permissionID", permissionID)
	return nil
}
