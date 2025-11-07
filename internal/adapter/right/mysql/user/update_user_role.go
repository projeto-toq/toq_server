package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserRole atualiza um user_role existente
func (ua *UserAdapter) UpdateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity, convertErr := userconverters.UserRoleDomainToEntity(userRole)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.user.update_user_role.convert_error", "error", convertErr)
		return fmt.Errorf("convert user role domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.permission.update_user_role.empty_entity")
		return nil
	}

	logger = logger.With(
		"user_role_id", entity.ID,
		"user_id", entity.UserID,
		"role_id", entity.RoleID,
	)

	query := `
		UPDATE user_roles 
		SET user_id = ?, role_id = ?, is_active = ?, status = ?, expires_at = ?
		WHERE id = ?
	`

	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		entity.UserID,
		entity.RoleID,
		entity.IsActive,
		entity.Status,
		entity.ExpiresAt,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.update_user_role.exec_error", "error", execErr)
		return fmt.Errorf("update user role: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.update_user_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user role update rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.update_user_role.success", "rows_affected", rowsAffected)
	return nil
}
