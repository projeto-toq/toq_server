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
func (p *PermissionAdapter) CreateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) (err error) {
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

	result, execErr := p.ExecContext(ctx, tx, "insert", query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.IsSystemRole,
		entity.IsActive,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.create_role.exec_error", "error", execErr)
		return fmt.Errorf("create role: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.permission.create_role.last_insert_id_error", "error", lastErr)
		return fmt.Errorf("role last insert id: %w", lastErr)
	}

	role.SetID(id)
	logger.Debug("mysql.permission.create_role.success", "role_id", id)
	return nil
}
