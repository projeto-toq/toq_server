package permissionconverters

import (
	"strings"

	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// PermissionDomainToEntity converte PermissionInterface para PermissionEntity
func PermissionDomainToEntity(permission permissionmodel.PermissionInterface) (*permissionentities.PermissionEntity, error) {
	if permission == nil {
		return nil, nil
	}

	entity := &permissionentities.PermissionEntity{
		ID:       permission.GetID(),
		Name:     permission.GetName(),
		Action:   permission.GetAction(),
		IsActive: permission.GetIsActive(),
	}

	if desc := strings.TrimSpace(permission.GetDescription()); desc != "" {
		entity.Description = &desc
	}

	return entity, nil
}
