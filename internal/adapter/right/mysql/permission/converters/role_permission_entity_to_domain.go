package permissionconverters

import (
	"encoding/json"
	"log/slog"

	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// RolePermissionEntityToDomain converte RolePermissionEntity para RolePermissionInterface
func RolePermissionEntityToDomain(entity *permissionentities.RolePermissionEntity) permissionmodel.RolePermissionInterface {
	if entity == nil {
		return nil
	}

	rolePermission := permissionmodel.NewRolePermission()
	rolePermission.SetID(entity.ID)
	rolePermission.SetRoleID(entity.RoleID)
	rolePermission.SetPermissionID(entity.PermissionID)
	rolePermission.SetGranted(entity.Granted)

	// Converter JSON conditions se existir
	if entity.Conditions != nil && *entity.Conditions != "" {
		var conditions map[string]interface{}
		if err := json.Unmarshal([]byte(*entity.Conditions), &conditions); err != nil {
			slog.Warn("Failed to unmarshal role permission conditions", "error", err, "conditions", *entity.Conditions)
		} else {
			rolePermission.SetConditions(conditions)
		}
	}

	return rolePermission
}
