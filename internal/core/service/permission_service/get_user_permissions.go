package permissionservice

import (
	"context"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetUserPermissions retorna todas as permissões de um usuário
func (p *permissionServiceImpl) GetUserPermissions(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	permissions, err := p.permissionRepository.GetUserPermissions(ctx, nil, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	return permissions, nil
}
