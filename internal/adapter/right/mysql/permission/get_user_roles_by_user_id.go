package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserRolesByUserID busca todos os user_roles de um usuário
func (pa *PermissionAdapter) GetUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (userRoles []permissionmodel.UserRoleInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID)

	query := `
		SELECT id, user_id, role_id, is_active, status, expires_at
		FROM user_roles 
		WHERE user_id = ?
		ORDER BY id
	`

	results, err := pa.Read(ctx, tx, query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_user_roles_by_user_id.read_error", "error", err)
		return nil, fmt.Errorf("get user roles by user id read: %w", err)
	}

	userRoles = make([]permissionmodel.UserRoleInterface, 0, len(results))
	for index, row := range results {
		if len(row) != 6 {
			errColumns := fmt.Errorf("unexpected number of columns: expected 6, got %d", len(row))
			utils.SetSpanError(ctx, errColumns)
			logger.Error("mysql.permission.get_user_roles_by_user_id.columns_mismatch", "row_index", index, "error", errColumns)
			return nil, errColumns
		}

		entity := &permissionentities.UserRoleEntity{
			ID:       row[0].(int64),
			UserID:   row[1].(int64),
			RoleID:   row[2].(int64),
			IsActive: row[3].(int64) == 1,
			Status:   row[4].(int64),
		}

		// Handle expires_at (pode ser NULL)
		if row[5] != nil {
			if expiresAt, ok := row[5].(time.Time); ok {
				entity.ExpiresAt = &expiresAt
			}
		}

		userRole, convertErr := permissionconverters.UserRoleEntityToDomain(entity)
		if convertErr != nil {
			utils.SetSpanError(ctx, convertErr)
			logger.Error("mysql.permission.get_user_roles_by_user_id.convert_error", "row_index", index, "error", convertErr)
			return nil, fmt.Errorf("convert user role entity to domain: %w", convertErr)
		}
		if userRole != nil {
			userRoles = append(userRoles, userRole)
		}
	}

	logger.Debug("mysql.permission.get_user_roles_by_user_id.success", "count", len(userRoles))
	return userRoles, nil
}
