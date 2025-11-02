package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdatePermission atualiza uma permissão existente
func (pa *PermissionAdapter) UpdatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity, convertErr := permissionconverters.PermissionDomainToEntity(permission)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.update_permission.convert_error", "error", convertErr)
		return fmt.Errorf("convert permission domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.permission.update_permission.empty_entity")
		return nil
	}

	logger = logger.With(
		"permission_id", entity.ID,
		"action", entity.Action,
	)

	query := `
		UPDATE permissions 
		SET name = ?, action = ?, description = ?, is_active = ?
		WHERE id = ?
	`

	result, execErr := pa.ExecContext(ctx, tx, "update", query,
		entity.Name,
		entity.Action,
		entity.Description,
		entity.IsActive,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.update_permission.exec_error", "error", execErr)
		return fmt.Errorf("update permission: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.update_permission.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("permission rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.update_permission.success", "rows_affected", rowsAffected)
	return nil
}
