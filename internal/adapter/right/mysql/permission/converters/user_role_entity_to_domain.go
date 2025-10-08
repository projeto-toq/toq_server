package permissionconverters

import (
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// UserRoleEntityToDomain converte UserRoleEntity para UserRoleInterface
func UserRoleEntityToDomain(entity *permissionentities.UserRoleEntity) (permissionmodel.UserRoleInterface, error) {
	if entity == nil {
		return nil, nil
	}

	userRole := permissionmodel.NewUserRole()
	userRole.SetID(entity.ID)
	userRole.SetUserID(entity.UserID)
	userRole.SetRoleID(entity.RoleID)
	userRole.SetIsActive(entity.IsActive)
	userRole.SetStatus(permissionmodel.UserRoleStatus(entity.Status))

	if entity.ExpiresAt != nil {
		userRole.SetExpiresAt(entity.ExpiresAt)
	}

	return userRole, nil
}
