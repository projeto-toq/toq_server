package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreatePermission cria uma nova permissão no banco de dados
func (pa *PermissionAdapter) CreatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity, convertErr := permissionconverters.PermissionDomainToEntity(permission)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.create_permission.convert_error", "error", convertErr)
		return fmt.Errorf("convert permission domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.permission.create_permission.empty_entity")
		return nil
	}

	logger = logger.With("action", entity.Action)

	query := `
		INSERT INTO permissions (name, action, description, is_active)
		VALUES (?, ?, ?, ?)
	`

	result, execErr := pa.ExecContext(ctx, tx, "insert", query,
		entity.Name,
		entity.Action,
		entity.Description,
		entity.IsActive,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.create_permission.exec_error", "error", execErr)
		return fmt.Errorf("create permission: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.permission.create_permission.last_insert_id_error", "error", lastErr)
		return fmt.Errorf("permission last insert id: %w", lastErr)
	}

	permission.SetID(id)
	logger.Debug("mysql.permission.create_permission.success", "permission_id", id)
	return nil
}
