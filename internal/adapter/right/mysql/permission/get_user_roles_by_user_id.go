package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
		SELECT 
			ur.id,
			ur.user_id,
			ur.role_id,
			ur.is_active,
			ur.status,
			ur.expires_at,
			r.id,
			r.slug,
			r.name,
			r.description,
			r.is_system_role,
			r.is_active
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ?
		ORDER BY ur.id
	`

	rows, readErr := pa.QueryContext(ctx, tx, "select", query, userID)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.permission.get_user_roles_by_user_id.read_error", "error", readErr)
		return nil, fmt.Errorf("get user roles by user id read: %w", readErr)
	}
	defer rows.Close()

	rowEntities, rowsErr := rowsToEntities(rows)
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.get_user_roles_by_user_id.rows_to_entities_error", "error", rowsErr)
		return nil, fmt.Errorf("get user roles by user id rows to entities: %w", rowsErr)
	}

	userRoles = make([]permissionmodel.UserRoleInterface, 0, len(rowEntities))
	for index, row := range rowEntities {
		if len(row) != 12 {
			errColumns := fmt.Errorf("unexpected number of columns: expected 12, got %d", len(row))
			utils.SetSpanError(ctx, errColumns)
			logger.Error("mysql.permission.get_user_roles_by_user_id.columns_mismatch", "row_index", index, "error", errColumns)
			return nil, errColumns
		}

		entity := &permissionentities.UserRoleEntity{
			ID:       int64FromAny(row[0]),
			UserID:   int64FromAny(row[1]),
			RoleID:   int64FromAny(row[2]),
			IsActive: intFromAny(row[3]) == 1,
			Status:   int64FromAny(row[4]),
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
			roleEntity := &permissionentities.RoleEntity{
				ID:           int64FromAny(row[6]),
				Slug:         stringFromAny(row[7]),
				Name:         stringFromAny(row[8]),
				Description:  stringFromNullable(row[9]),
				IsSystemRole: intFromAny(row[10]) == 1,
				IsActive:     intFromAny(row[11]) == 1,
			}
			roleDomain := permissionconverters.RoleEntityToDomain(roleEntity)
			if roleDomain != nil {
				userRole.SetRole(roleDomain)
			}
			userRoles = append(userRoles, userRole)
		}
	}

	logger.Debug("mysql.permission.get_user_roles_by_user_id.success", "count", len(userRoles))
	return userRoles, nil
}

func stringFromAny(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func stringFromNullable(value interface{}) string {
	return stringFromAny(value)
}

func intFromAny(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case []byte:
		if i, err := strconv.Atoi(string(v)); err == nil {
			return i
		}
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

func int64FromAny(value interface{}) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case uint64:
		return int64(v)
	case uint32:
		return int64(v)
	case uint:
		return int64(v)
	case []byte:
		if i, err := strconv.ParseInt(string(v), 10, 64); err == nil {
			return i
		}
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}
