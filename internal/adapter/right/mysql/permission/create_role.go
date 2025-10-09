package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateRole cria um novo role no banco de dados
func (pa *PermissionAdapter) CreateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity := permissionconverters.RoleDomainToEntity(role)
	if entity == nil {
		logger.Warn("mysql.permission.create_role.empty_entity")
		return nil
	}

	logger = logger.With(
		"role_slug", entity.Slug,
		"role_name", entity.Name,
	)

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
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.create_role.exec_error", "error", err)
		return fmt.Errorf("create role: %w", err)
	}

	role.SetID(id)
	logger.Debug("mysql.permission.create_role.success", "role_id", id)
	return nil
}
