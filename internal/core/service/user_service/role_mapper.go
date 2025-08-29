package userservices

import (
	"context"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// getRoleBySlug busca um role pelo slug usando o permission service
func (us *userService) getRoleBySlug(ctx context.Context, slug string) (permissionmodel.RoleInterface, error) {
	if us.permissionService == nil {
		return nil, fmt.Errorf("permission service not available")
	}

	role, err := us.permissionService.GetRoleBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by slug '%s': %w", slug, err)
	}

	return role, nil
}

// getLegacyRoleBySlug converte slug para usermodel.UserRole (para compatibilidade com métodos antigos)
func (us *userService) getLegacyRoleBySlug(roleSlug string) (usermodel.UserRole, error) {
	switch roleSlug {
	case "admin":
		return usermodel.RoleRoot, nil
	case "owner":
		return usermodel.RoleOwner, nil
	case "realtor":
		return usermodel.RoleRealtor, nil
	case "agency":
		return usermodel.RoleAgency, nil
	default:
		return usermodel.UserRole(0), fmt.Errorf("unknown role slug: %s", roleSlug)
	}
}

// assignRoleToUser atribui um role ao usuário usando o permission service
func (us *userService) assignRoleToUser(ctx context.Context, userID int64, roleSlug string) error {
	if us.permissionService == nil {
		return fmt.Errorf("permission service not available")
	}

	role, err := us.getRoleBySlug(ctx, roleSlug)
	if err != nil {
		return err
	}

	err = us.permissionService.AssignRoleToUser(ctx, userID, role.GetID(), nil)
	if err != nil {
		return fmt.Errorf("failed to assign role '%s' to user %d: %w", roleSlug, userID, err)
	}

	return nil
}
