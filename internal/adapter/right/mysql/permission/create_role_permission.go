package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// CreateRolePermission cria uma nova associação role-permission no banco de dados
func (pa *PermissionAdapter) CreateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) error {
	entity := permissionconverters.RolePermissionDomainToEntity(rolePermission)
	if entity == nil {
		return nil
	}

	query := `
		INSERT INTO role_permissions (role_id, permission_id, granted, conditions)
		VALUES (?, ?, ?, ?)
	`

	id, err := pa.Create(ctx, tx, query,
		entity.RoleID,
		entity.PermissionID,
		entity.Granted,
		entity.Conditions,
	)
	if err != nil {
		return err
	}

	rolePermission.SetID(id)
	return nil
}
