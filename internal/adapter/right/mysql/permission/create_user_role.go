package mysqlpermissionadapter

import (
	"context"
	"database/sql"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// CreateUserRole cria uma nova associação user-role no banco de dados
func (pa *PermissionAdapter) CreateUserRole(ctx context.Context, tx *sql.Tx, userRole permissionmodel.UserRoleInterface) error {
	entity := permissionconverters.UserRoleDomainToEntity(userRole)
	if entity == nil {
		return nil
	}

	query := `
		INSERT INTO user_roles (user_id, role_id, is_active, expires_at)
		VALUES (?, ?, ?, ?)
	`

	id, err := pa.Create(ctx, tx, query,
		entity.UserID,
		entity.RoleID,
		entity.IsActive,
		entity.ExpiresAt,
	)
	if err != nil {
		return err
	}

	userRole.SetID(id)
	return nil
}
