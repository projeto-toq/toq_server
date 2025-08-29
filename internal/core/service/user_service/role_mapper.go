package userservices

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// getRoleBySlug busca um role pelo slug usando o permission service
func (us *userService) getRoleBySlug(ctx context.Context, tx *sql.Tx, slug string) (permissionmodel.RoleInterface, error) {
	if us.permissionService == nil {
		return nil, utils.ErrInternalServer
	}

	role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, slug)
	if err != nil {
		return nil, utils.ErrInternalServer
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
		return usermodel.UserRole(0), utils.ErrBadRequest
	}
}

// assignRoleToUser atribui um role ao usuário usando o permission service
func (us *userService) assignRoleToUser(ctx context.Context, tx *sql.Tx, userID int64, roleSlug string) error {
	if us.permissionService == nil {
		return utils.ErrInternalServer
	}

	role, err := us.getRoleBySlug(ctx, tx, roleSlug)
	if err != nil {
		return err
	}

	err = us.permissionService.AssignRoleToUserWithTx(ctx, tx, userID, role.GetID(), nil)
	if err != nil {
		return utils.ErrInternalServer
	}

	return nil
}
