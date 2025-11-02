package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateRolePermission atualiza um role_permission existente
func (pa *PermissionAdapter) UpdateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity, convertErr := permissionconverters.RolePermissionDomainToEntity(rolePermission)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.update_role_permission.convert_error", "error", convertErr)
		return fmt.Errorf("convert role permission domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.permission.update_role_permission.empty_entity")
		return nil
	}

	logger = logger.With(
		"role_permission_id", entity.ID,
		"role_id", entity.RoleID,
		"permission_id", entity.PermissionID,
	)

	query := `
		UPDATE role_permissions 
		SET role_id = ?, permission_id = ?, granted = ?
		WHERE id = ?
	`

	result, execErr := pa.ExecContext(ctx, tx, "update", query,
		entity.RoleID,
		entity.PermissionID,
		entity.Granted,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.update_role_permission.exec_error", "error", execErr)
		return fmt.Errorf("update role permission: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.update_role_permission.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("role permission update rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.update_role_permission.success", "rows_affected", rowsAffected)
	return nil
}
