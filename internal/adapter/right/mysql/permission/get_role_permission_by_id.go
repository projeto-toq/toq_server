package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRolePermissionByID retorna uma relação role-permission pelo ID
func (pa *PermissionAdapter) GetRolePermissionByID(ctx context.Context, tx *sql.Tx, rolePermissionID int64) (permissionmodel.RolePermissionInterface, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("role_permission_id", rolePermissionID)

	query := `SELECT id, role_id, permission_id, granted FROM role_permissions WHERE id = ?`

	rows, readErr := pa.QueryContext(ctx, tx, "select", query, rolePermissionID)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.permission.get_role_permission_by_id.read_error", "error", readErr)
		return nil, fmt.Errorf("get role permission by id read: %w", readErr)
	}
	defer rows.Close()

	rowEntities, rowsErr := rowsToEntities(rows)
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.get_role_permission_by_id.rows_to_entities_error", "error", rowsErr)
		return nil, fmt.Errorf("get role permission by id rows to entities: %w", rowsErr)
	}
	if len(rowEntities) == 0 {
		return nil, nil
	}

	row := rowEntities[0]
	if len(row) != 4 {
		logger.Warn("mysql.permission.get_role_permission_by_id.columns_mismatch", "expected", 4, "got", len(row))
		return nil, fmt.Errorf("unexpected number of columns")
	}

	entity := &permissionentities.RolePermissionEntity{}
	if val, ok := row[0].(int64); ok {
		entity.ID = val
	}
	if val, ok := row[1].(int64); ok {
		entity.RoleID = val
	}
	if val, ok := row[2].(int64); ok {
		entity.PermissionID = val
	}
	switch grantedVal := row[3].(type) {
	case int64:
		entity.Granted = grantedVal == 1
	case bool:
		entity.Granted = grantedVal
	}

	rolePermission, convertErr := permissionconverters.RolePermissionEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_role_permission_by_id.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert role permission entity to domain: %w", convertErr)
	}

	return rolePermission, nil
}
