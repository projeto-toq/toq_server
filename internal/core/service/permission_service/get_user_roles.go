package permissionservice

import (
	"context"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetUserRoles retorna todos os roles ativos de um usuário
func (p *permissionServiceImpl) GetUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	userRoles, err := p.permissionRepository.GetActiveUserRolesByUserID(ctx, nil, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return userRoles, nil
}
