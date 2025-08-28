package permissionconverters

import (
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// RoleDomainToEntity converte RoleInterface para RoleEntity
func RoleDomainToEntity(role permissionmodel.RoleInterface) *permissionentities.RoleEntity {
	if role == nil {
		return nil
	}

	return &permissionentities.RoleEntity{
		ID:           role.GetID(),
		Name:         role.GetName(),
		Slug:         role.GetSlug(),
		Description:  role.GetDescription(),
		IsSystemRole: role.GetIsSystemRole(),
		IsActive:     role.GetIsActive(),
	}
}
