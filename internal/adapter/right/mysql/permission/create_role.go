package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// CreateRole cria um novo role no banco de dados
func (pa *PermissionAdapter) CreateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) error {
	entity := permissionconverters.RoleDomainToEntity(role)
	if entity == nil {
		return nil
	}

	query := `
		INSERT INTO roles (name, slug, description, is_system_role, is_active)
		VALUES (?, ?, ?, ?, ?)
	`

	id, err := pa.Create(ctx, tx, query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.IsSystemRole,
		entity.IsActive,
	)
	if err != nil {
		return err
	}

	role.SetID(id)
	return nil
}
