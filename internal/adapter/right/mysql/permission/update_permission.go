package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
		"resource", entity.Resource,
		"action", entity.Action,
	)

	query := `
		UPDATE permissions 
		SET name = ?, resource = ?, action = ?, description = ?, conditions = ?, is_active = ?
		WHERE id = ?
	`

	rowsAffected, err := pa.Update(ctx, tx, query,
		entity.Name,
		entity.Resource,
		entity.Action,
		entity.Description,
		entity.Conditions,
		entity.IsActive,
		entity.ID,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.update_permission.exec_error", "error", err)
		return fmt.Errorf("update permission: %w", err)
	}

	logger.Debug("mysql.permission.update_permission.success", "rows_affected", rowsAffected)
	return nil
}
