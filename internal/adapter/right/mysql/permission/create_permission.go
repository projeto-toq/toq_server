package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// CreatePermission cria uma nova permissão no banco de dados
func (pa *PermissionAdapter) CreatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) error {
	entity := permissionconverters.PermissionDomainToEntity(permission)
	if entity == nil {
		return nil
	}

	query := `
		INSERT INTO permissions (name, slug, resource, action, description, conditions, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	id, err := pa.Create(ctx, tx, query,
		entity.Name,
		entity.Slug,
		entity.Resource,
		entity.Action,
		entity.Description,
		entity.Conditions,
		entity.IsActive,
	)
	if err != nil {
		return err
	}

	permission.SetID(id)
	return nil
}
