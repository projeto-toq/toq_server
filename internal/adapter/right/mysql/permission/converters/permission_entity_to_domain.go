package permissionconverters

import (
	"encoding/json"
	"fmt"

	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// PermissionEntityToDomain converte PermissionEntity para PermissionInterface
func PermissionEntityToDomain(entity *permissionentities.PermissionEntity) (permissionmodel.PermissionInterface, error) {
	if entity == nil {
		return nil, nil
	}

	permission := permissionmodel.NewPermission()
	permission.SetID(entity.ID)
	permission.SetName(entity.Name)
	permission.SetResource(entity.Resource)
	permission.SetAction(entity.Action)
	permission.SetDescription(entity.Description)

	// Converter JSON conditions se existir
	if entity.Conditions != nil && *entity.Conditions != "" {
		var conditions map[string]interface{}
		if err := json.Unmarshal([]byte(*entity.Conditions), &conditions); err != nil {
			return permission, fmt.Errorf("unmarshal permission conditions: %w", err)
		}
		permission.SetConditions(conditions)
	}

	return permission, nil
}
