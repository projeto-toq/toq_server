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

// GetUserRoleByUserIDAndRoleID busca um user_role específico pela combinação user_id + role_id
func (pa *PermissionAdapter) GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (permissionmodel.UserRoleInterface, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID, "role_id", roleID)

	query := `
		SELECT id, user_id, role_id, is_active, status, expires_at
		FROM user_roles 
		WHERE user_id = ? AND role_id = ?
		LIMIT 1
	`

	var (
		id          int64
		uid         int64
		roleIDOut   int64
		isActiveInt int64
		status      int64
		expiresAt   sql.NullTime
	)

	row := pa.QueryRowContext(ctx, tx, "select", query, userID, roleID)
	err = row.Scan(
		&id, &uid, &roleIDOut, &isActiveInt, &status, &expiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_user_role_by_user_id_and_role_id.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_user_role_by_user_id_and_role_id.scan_error", "error", err)
		return nil, fmt.Errorf("get user role by user id and role id scan: %w", err)
	}

	entity := &permissionentities.UserRoleEntity{
		ID:       id,
		UserID:   uid,
		RoleID:   roleIDOut,
		IsActive: isActiveInt == 1,
		Status:   status,
	}
	if expiresAt.Valid {
		t := expiresAt.Time
		entity.ExpiresAt = &t
	}

	userRole, convertErr := permissionconverters.UserRoleEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_user_role_by_user_id_and_role_id.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert user role entity to domain: %w", convertErr)
	}

	logger.Debug("mysql.permission.get_user_role_by_user_id_and_role_id.success")
	return userRole, nil
}
