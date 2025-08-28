package permissionconverters

import (
	"encoding/json"
	"log/slog"

	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// RolePermissionDomainToEntity converte RolePermissionInterface para RolePermissionEntity
func RolePermissionDomainToEntity(rolePermission permissionmodel.RolePermissionInterface) *permissionentities.RolePermissionEntity {
	if rolePermission == nil {
		return nil
	}

	entity := &permissionentities.RolePermissionEntity{
		ID:           rolePermission.GetID(),
		RoleID:       rolePermission.GetRoleID(),
		PermissionID: rolePermission.GetPermissionID(),
		Granted:      rolePermission.GetGranted(),
	}

	// Converter conditions para JSON se existir
	if conditions := rolePermission.GetConditions(); conditions != nil {
		conditionsJSON, err := json.Marshal(conditions)
		if err != nil {
			slog.Warn("Failed to marshal role permission conditions", "error", err)
		} else {
			conditionsStr := string(conditionsJSON)
			entity.Conditions = &conditionsStr
		}
	}

	return entity
}
