package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetRolePermissionsByRoleID busca todas as associações role_permission de um role
func (pa *PermissionAdapter) GetRolePermissionsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) (rolePermissions []permissionmodel.RolePermissionInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("role_id", roleID)

	query := `
		SELECT id, role_id, permission_id, granted, conditions
		FROM role_permissions 
		WHERE role_id = ?
		ORDER BY id
	`

	results, err := pa.Read(ctx, tx, query, roleID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_role_permissions_by_role_id.read_error", "error", err)
		return nil, fmt.Errorf("get role permissions by role id read: %w", err)
	}

	rolePermissions = make([]permissionmodel.RolePermissionInterface, 0, len(results))
	for index, row := range results {
		if len(row) != 5 {
			errColumns := fmt.Errorf("unexpected number of columns: expected 5, got %d", len(row))
			utils.SetSpanError(ctx, errColumns)
			logger.Error("mysql.permission.get_role_permissions_by_role_id.columns_mismatch", "row_index", index, "error", errColumns)
			return nil, errColumns
		}

		entity := &permissionentities.RolePermissionEntity{
			ID:           row[0].(int64),
			RoleID:       row[1].(int64),
			PermissionID: row[2].(int64),
			Granted:      row[3].(int64) == 1,
		}

		// conditions (pode ser NULL)
		if row[4] != nil {
			conditionsStr := string(row[4].([]byte))
			entity.Conditions = &conditionsStr
		}

		rolePermission, convertErr := permissionconverters.RolePermissionEntityToDomain(entity)
		if convertErr != nil {
			utils.SetSpanError(ctx, convertErr)
			logger.Error("mysql.permission.get_role_permissions_by_role_id.convert_error", "row_index", index, "error", convertErr)
			return nil, fmt.Errorf("convert role permission entity to domain: %w", convertErr)
		}
		if rolePermission != nil {
			rolePermissions = append(rolePermissions, rolePermission)
		}
	}

	logger.Debug("mysql.permission.get_role_permissions_by_role_id.success", "count", len(rolePermissions))
	return rolePermissions, nil
}
