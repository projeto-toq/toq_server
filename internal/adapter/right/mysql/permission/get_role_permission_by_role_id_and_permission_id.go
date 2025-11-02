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

// GetRolePermissionByRoleIDAndPermissionID busca um role_permission específico pela combinação role_id + permission_id
func (pa *PermissionAdapter) GetRolePermissionByRoleIDAndPermissionID(ctx context.Context, tx *sql.Tx, roleID, permissionID int64) (rolePermission permissionmodel.RolePermissionInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("role_id", roleID, "permission_id", permissionID)

	query := `
		SELECT id, role_id, permission_id, granted
		FROM role_permissions 
		WHERE role_id = ? AND permission_id = ?
		LIMIT 1
	`

	var (
		id              int64
		roleIDOut       int64
		permissionIDOut int64
		grantedInt      int64
	)

	row := pa.QueryRowContext(ctx, tx, "select", query, roleID, permissionID)
	err = row.Scan(
		&id, &roleIDOut, &permissionIDOut, &grantedInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_role_permission_by_ids.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_role_permission_by_ids.scan_error", "error", err)
		return nil, fmt.Errorf("get role permission by ids scan: %w", err)
	}

	entity := &permissionentities.RolePermissionEntity{
		ID:           id,
		RoleID:       roleIDOut,
		PermissionID: permissionIDOut,
		Granted:      grantedInt == 1,
	}

	rolePermission, convertErr := permissionconverters.RolePermissionEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_role_permission_by_ids.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert role permission entity to domain: %w", convertErr)
	}

	logger.Debug("mysql.permission.get_role_permission_by_ids.success")
	return rolePermission, nil
}
