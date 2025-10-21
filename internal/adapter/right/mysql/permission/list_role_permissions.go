package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListRolePermissions retorna relações role-permission aplicando filtros e paginação
func (pa *PermissionAdapter) ListRolePermissions(ctx context.Context, tx *sql.Tx, filter permissionrepository.RolePermissionListFilter) (permissionrepository.RolePermissionListResult, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return permissionrepository.RolePermissionListResult{}, err
	}
	defer spanEnd()

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	conditions := []string{"1=1"}
	args := []any{}

	if filter.RoleID != nil {
		conditions = append(conditions, "rp.role_id = ?")
		args = append(args, *filter.RoleID)
	}
	if filter.PermissionID != nil {
		conditions = append(conditions, "rp.permission_id = ?")
		args = append(args, *filter.PermissionID)
	}
	if filter.Granted != nil {
		if *filter.Granted {
			conditions = append(conditions, "rp.granted = 1")
		} else {
			conditions = append(conditions, "rp.granted = 0")
		}
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	query := `SELECT id, role_id, permission_id, granted
        FROM role_permissions rp ` + " " + whereClause + " ORDER BY id ASC LIMIT ? OFFSET ?"

	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	rows, readErr := pa.Read(ctx, tx, query, listArgs...)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.permission.list_role_permissions.read_error", "error", readErr)
		return permissionrepository.RolePermissionListResult{}, fmt.Errorf("list role permissions read: %w", readErr)
	}

	result := permissionrepository.RolePermissionListResult{}
	if len(rows) > 0 {
		result.RolePermissions = make([]permissionmodel.RolePermissionInterface, 0, len(rows))
	}

	for idx, row := range rows {
		if len(row) != 4 {
			logger.Warn("mysql.permission.list_role_permissions.columns_mismatch", "row_index", idx, "expected", 4, "got", len(row))
			continue
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
			logger.Error("mysql.permission.list_role_permissions.convert_error", "row_index", idx, "error", convertErr)
			return permissionrepository.RolePermissionListResult{}, fmt.Errorf("convert role permission entity to domain: %w", convertErr)
		}
		if rolePermission != nil {
			result.RolePermissions = append(result.RolePermissions, rolePermission)
		}
	}

	countQuery := `SELECT COUNT(*) FROM role_permissions rp ` + whereClause
	countRows, countErr := pa.Read(ctx, tx, countQuery, args...)
	if countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.permission.list_role_permissions.count_error", "error", countErr)
		return permissionrepository.RolePermissionListResult{}, fmt.Errorf("list role permissions count: %w", countErr)
	}
	if len(countRows) > 0 && len(countRows[0]) > 0 {
		switch total := countRows[0][0].(type) {
		case int64:
			result.Total = total
		case []byte:
			if val, convErr := strconv.ParseInt(string(total), 10, 64); convErr == nil {
				result.Total = val
			}
		}
	}

	return result, nil
}
