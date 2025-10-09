package permissionconverters

import (
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// RoleEntityToDomain converte RoleEntity para RoleInterface
func RoleEntityToDomain(entity *permissionentities.RoleEntity) permissionmodel.RoleInterface {
	if entity == nil {
		return nil
	}

	role := permissionmodel.NewRole()
	role.SetID(entity.ID)
	role.SetName(entity.Name)
	role.SetSlug(entity.Slug)
	role.SetDescription(entity.Description)
	role.SetIsSystemRole(entity.IsSystemRole)
	role.SetIsActive(entity.IsActive)

	return role
}
