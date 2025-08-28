package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// UpdateRolePermission atualiza um role_permission existente
func (pa *PermissionAdapter) UpdateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) error {
	entity := permissionconverters.RolePermissionDomainToEntity(rolePermission)
	if entity == nil {
		return nil
	}

	query := `
		UPDATE role_permissions 
		SET role_id = ?, permission_id = ?, granted = ?, conditions = ?
		WHERE id = ?
	`

	_, err := pa.Update(ctx, tx, query,
		entity.RoleID,
		entity.PermissionID,
		entity.Granted,
		entity.Conditions,
		entity.ID,
	)

	return err
}
