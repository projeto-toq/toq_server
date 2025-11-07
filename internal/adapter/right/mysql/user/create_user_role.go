package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateUserRole cria uma nova associação user-role no banco de dados
func (ua *UserAdapter) CreateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) (result usermodel.UserRoleInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity, convertErr := userconverters.UserRoleDomainToEntity(userRole)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.user.create_user_role.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert user role domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.user.create_user_role.empty_entity")
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

	resultExec, execErr := ua.ExecContext(ctx, tx, "insert", query,
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
