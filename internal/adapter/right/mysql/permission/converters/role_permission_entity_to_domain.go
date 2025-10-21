package permissionconverters

import (
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// RolePermissionEntityToDomain converte RolePermissionEntity para RolePermissionInterface
func RolePermissionEntityToDomain(entity *permissionentities.RolePermissionEntity) (permissionmodel.RolePermissionInterface, error) {
	if entity == nil {
		return nil, nil
	}

	rolePermission := permissionmodel.NewRolePermission()
	rolePermission.SetID(entity.ID)
	rolePermission.SetRoleID(entity.RoleID)
	rolePermission.SetPermissionID(entity.PermissionID)
	rolePermission.SetGranted(entity.Granted)

	return rolePermission, nil
}
