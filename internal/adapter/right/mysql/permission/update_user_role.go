package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// UpdateUserRole atualiza um user_role existente
func (pa *PermissionAdapter) UpdateUserRole(ctx context.Context, tx *sql.Tx, userRole permissionmodel.UserRoleInterface) error {
	entity := permissionconverters.UserRoleDomainToEntity(userRole)
	if entity == nil {
		return nil
	}

	query := `
		UPDATE user_roles 
		SET user_id = ?, role_id = ?, is_active = ?, expires_at = ?
		WHERE id = ?
	`

	_, err := pa.Update(ctx, tx, query,
		entity.UserID,
		entity.RoleID,
		entity.IsActive,
		entity.ExpiresAt,
		entity.ID,
	)

	return err
}
