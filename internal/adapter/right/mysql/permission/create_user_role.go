package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateUserRole cria uma nova associação user-role no banco de dados
func (pa *PermissionAdapter) CreateUserRole(ctx context.Context, tx *sql.Tx, userRole permissionmodel.UserRoleInterface) (result permissionmodel.UserRoleInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	entity, convertErr := permissionconverters.UserRoleDomainToEntity(userRole)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.create_user_role.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert user role domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.permission.create_user_role.empty_entity")
		return nil, nil
	}

	logger = logger.With(
		"user_id", entity.UserID,
		"role_id", entity.RoleID,
	)

	query := `
		INSERT INTO user_roles (user_id, role_id, is_active, status, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`

	resultExec, execErr := pa.ExecContext(ctx, tx, "insert", query,
		entity.UserID,
		entity.RoleID,
		entity.IsActive,
		entity.Status,
		entity.ExpiresAt,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.create_user_role.exec_error", "error", execErr)
		return nil, fmt.Errorf("create user role: %w", execErr)
	}

	id, lastErr := resultExec.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.permission.create_user_role.last_insert_id_error", "error", lastErr)
		return nil, fmt.Errorf("user role last insert id: %w", lastErr)
	}

	userRole.SetID(id)
	logger.Debug("mysql.permission.create_user_role.success", "user_role_id", id)

	return userRole, nil
}
