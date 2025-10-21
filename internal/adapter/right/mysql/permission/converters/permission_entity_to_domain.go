package permissionconverters

import (
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// PermissionEntityToDomain converte PermissionEntity para PermissionInterface
func PermissionEntityToDomain(entity *permissionentities.PermissionEntity) (permissionmodel.PermissionInterface, error) {
	if entity == nil {
		return nil, nil
	}

	permission := permissionmodel.NewPermission()
	permission.SetID(entity.ID)
	permission.SetName(entity.Name)
	permission.SetAction(entity.Action)
	if entity.Description != nil {
		permission.SetDescription(*entity.Description)
	}
	permission.SetIsActive(entity.IsActive)

	return permission, nil
}
