package permissionconverters

import (
	"encoding/json"
	"log/slog"

	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// PermissionEntityToDomain converte PermissionEntity para PermissionInterface
func PermissionEntityToDomain(entity *permissionentities.PermissionEntity) permissionmodel.PermissionInterface {
	if entity == nil {
		return nil
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
			slog.Warn("Failed to unmarshal permission conditions", "error", err, "conditions", *entity.Conditions)
		} else {
			permission.SetConditions(conditions)
		}
	}

	return permission
}
