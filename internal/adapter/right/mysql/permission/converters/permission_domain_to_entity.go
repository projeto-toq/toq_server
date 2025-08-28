package permissionconverters

import (
	"encoding/json"
	"log/slog"

	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// PermissionDomainToEntity converte PermissionInterface para PermissionEntity
func PermissionDomainToEntity(permission permissionmodel.PermissionInterface) *permissionentities.PermissionEntity {
	if permission == nil {
		return nil
	}

	entity := &permissionentities.PermissionEntity{
		ID:          permission.GetID(),
		Name:        permission.GetName(),
		Slug:        "", // TODO: Implementar slug se necessário
		Resource:    permission.GetResource(),
		Action:      permission.GetAction(),
		Description: permission.GetDescription(),
		IsActive:    true, // Default to active
	}

	// Converter conditions para JSON se existir
	if conditions := permission.GetConditions(); conditions != nil {
		conditionsJSON, err := json.Marshal(conditions)
		if err != nil {
			slog.Warn("Failed to marshal permission conditions", "error", err)
		} else {
			conditionsStr := string(conditionsJSON)
			entity.Conditions = &conditionsStr
		}
	}

	return entity
}
