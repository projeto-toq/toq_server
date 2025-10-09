package permissionconverters

import (
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// UserRoleDomainToEntity converte UserRoleInterface para UserRoleEntity
func UserRoleDomainToEntity(userRole permissionmodel.UserRoleInterface) (*permissionentities.UserRoleEntity, error) {
	if userRole == nil {
		return nil, nil
	}

	entity := &permissionentities.UserRoleEntity{
		ID:       userRole.GetID(),
		UserID:   userRole.GetUserID(),
		RoleID:   userRole.GetRoleID(),
		IsActive: userRole.GetIsActive(),
		Status:   int64(userRole.GetStatus()),
	}

	if expiresAt := userRole.GetExpiresAt(); expiresAt != nil {
		entity.ExpiresAt = expiresAt
	}

	return entity, nil
}
