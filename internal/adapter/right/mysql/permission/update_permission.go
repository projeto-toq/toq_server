package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// UpdatePermission atualiza uma permissão existente
func (pa *PermissionAdapter) UpdatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) error {
	entity := permissionconverters.PermissionDomainToEntity(permission)
	if entity == nil {
		return nil
	}

	query := `
		UPDATE permissions 
		SET name = ?, resource = ?, action = ?, description = ?, conditions = ?, is_active = ?
		WHERE id = ?
	`

	_, err := pa.Update(ctx, tx, query,
		entity.Name,
		entity.Resource,
		entity.Action,
		entity.Description,
		entity.Conditions,
		entity.IsActive,
		entity.ID,
	)

	return err
}
