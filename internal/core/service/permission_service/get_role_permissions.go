package permissionservice

import (
	"context"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetRolePermissions retorna todas as permissões de um role
func (p *permissionServiceImpl) GetRolePermissions(ctx context.Context, roleID int64) ([]permissionmodel.PermissionInterface, error) {
	if roleID <= 0 {
		return nil, fmt.Errorf("invalid role ID: %d", roleID)
	}

	permissions, err := p.permissionRepository.GetGrantedPermissionsByRoleID(ctx, nil, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return permissions, nil
}
