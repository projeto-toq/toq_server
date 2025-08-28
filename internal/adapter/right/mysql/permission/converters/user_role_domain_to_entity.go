package permissionconverters

import (
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// UserRoleDomainToEntity converte UserRoleInterface para UserRoleEntity
func UserRoleDomainToEntity(userRole permissionmodel.UserRoleInterface) *permissionentities.UserRoleEntity {
	if userRole == nil {
		return nil
	}

	entity := &permissionentities.UserRoleEntity{
		ID:       userRole.GetID(),
		UserID:   userRole.GetUserID(),
		RoleID:   userRole.GetRoleID(),
		IsActive: userRole.GetIsActive(),
	}

	if expiresAt := userRole.GetExpiresAt(); expiresAt != nil {
		entity.ExpiresAt = expiresAt
	}

	return entity
}
