package permissionconverters

import (
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// RolePermissionDomainToEntity converte RolePermissionInterface para RolePermissionEntity
func RolePermissionDomainToEntity(rolePermission permissionmodel.RolePermissionInterface) (*permissionentities.RolePermissionEntity, error) {
	if rolePermission == nil {
		return nil, nil
	}

	entity := &permissionentities.RolePermissionEntity{
		ID:           rolePermission.GetID(),
		RoleID:       rolePermission.GetRoleID(),
		PermissionID: rolePermission.GetPermissionID(),
		Granted:      rolePermission.GetGranted(),
	}

	return entity, nil
}
