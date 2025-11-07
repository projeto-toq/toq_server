package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateRolePermission cria uma nova associação role-permission no banco de dados
func (p *PermissionAdapter) CreateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity, convertErr := permissionconverters.RolePermissionDomainToEntity(rolePermission)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.create_role_permission.convert_error", "error", convertErr)
		return fmt.Errorf("convert role permission domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.permission.create_role_permission.empty_entity")
		return nil
	}

	logger = logger.With(
		"role_id", entity.RoleID,
		"permission_id", entity.PermissionID,
	)

	query := `
		INSERT INTO role_permissions (role_id, permission_id, granted)
		VALUES (?, ?, ?)
	`

	result, execErr := p.ExecContext(ctx, tx, "insert", query,
		entity.RoleID,
		entity.PermissionID,
		entity.Granted,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.create_role_permission.exec_error", "error", execErr)
		return fmt.Errorf("create role permission: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.permission.create_role_permission.last_insert_id_error", "error", lastErr)
		return fmt.Errorf("role permission last insert id: %w", lastErr)
	}

	rolePermission.SetID(id)
	logger.Debug("mysql.permission.create_role_permission.success", "role_permission_id", id)
	return nil
}
