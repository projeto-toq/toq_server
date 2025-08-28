package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// UpdateRole atualiza um role existente
func (pa *PermissionAdapter) UpdateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) error {
	entity := permissionconverters.RoleDomainToEntity(role)
	if entity == nil {
		return nil
	}

	query := `
		UPDATE roles 
		SET name = ?, slug = ?, description = ?, is_system_role = ?, is_active = ?
		WHERE id = ?
	`

	_, err := pa.Update(ctx, tx, query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.IsSystemRole,
		entity.IsActive,
		entity.ID,
	)

	return err
}
