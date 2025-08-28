package permissionservice

import (
	"context"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GrantPermissionToRole concede uma permissão a um role
func (p *permissionServiceImpl) GrantPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	if roleID <= 0 {
		return fmt.Errorf("invalid role ID: %d", roleID)
	}

	if permissionID <= 0 {
		return fmt.Errorf("invalid permission ID: %d", permissionID)
	}

	slog.Debug("Granting permission to role", "roleID", roleID, "permissionID", permissionID)

	// Verificar se o role existe
	role, err := p.permissionRepository.GetRoleByID(ctx, nil, roleID)
	if err != nil {
		return fmt.Errorf("failed to verify role existence: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role with ID %d does not exist", roleID)
	}

	// Verificar se a permissão existe
	permission, err := p.permissionRepository.GetPermissionByID(ctx, nil, permissionID)
	if err != nil {
		return fmt.Errorf("failed to verify permission existence: %w", err)
	}
	if permission == nil {
		return fmt.Errorf("permission with ID %d does not exist", permissionID)
	}

	// Verificar se a relação já existe
	existingRolePermission, err := p.permissionRepository.GetRolePermissionByRoleIDAndPermissionID(ctx, nil, roleID, permissionID)
	if err == nil && existingRolePermission != nil {
		return fmt.Errorf("role %d already has permission %d", roleID, permissionID)
	}

	// Criar a nova RolePermission
	rolePermission := permissionmodel.NewRolePermission()
	rolePermission.SetRoleID(roleID)
	rolePermission.SetPermissionID(permissionID)
	rolePermission.SetGranted(true)

	// Salvar no banco
	err = p.permissionRepository.CreateRolePermission(ctx, nil, rolePermission)
	if err != nil {
		return fmt.Errorf("failed to grant permission to role: %w", err)
	}

	slog.Info("Permission granted to role successfully", "roleID", roleID, "permissionID", permissionID)
	return nil
}
