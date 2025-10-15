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

// ListPermissions retorna permissões com filtros e paginação
func (pa *PermissionAdapter) ListPermissions(ctx context.Context, tx *sql.Tx, filter permissionrepository.PermissionListFilter) (permissionrepository.PermissionListResult, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return permissionrepository.PermissionListResult{}, err
	}
	defer spanEnd()

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	conditions := []string{"1=1"}
	var args []any

	if filter.Name != "" {
		conditions = append(conditions, "p.name LIKE ?")
		args = append(args, filter.Name)
	}
	if filter.Resource != "" {
		conditions = append(conditions, "p.resource LIKE ?")
		args = append(args, filter.Resource)
	}
	if filter.Action != "" {
		conditions = append(conditions, "p.action LIKE ?")
		args = append(args, filter.Action)
	}
	if filter.IsActive != nil {
		if *filter.IsActive {
			conditions = append(conditions, "p.is_active = 1")
		} else {
			conditions = append(conditions, "p.is_active = 0")
		}
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	query := `SELECT id, name, resource, action, description, conditions, is_active
        FROM permissions p ` + " " + whereClause + " ORDER BY resource ASC, action ASC LIMIT ? OFFSET ?"

	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	rows, readErr := pa.Read(ctx, tx, query, listArgs...)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.permission.list_permissions.read_error", "error", readErr)
		return permissionrepository.PermissionListResult{}, fmt.Errorf("list permissions read: %w", readErr)
	}

	result := permissionrepository.PermissionListResult{}
	if len(rows) > 0 {
		result.Permissions = make([]permissionmodel.PermissionInterface, 0, len(rows))
	}

	for idx, row := range rows {
		if len(row) != 7 {
			logger.Warn("mysql.permission.list_permissions.columns_mismatch", "row_index", idx, "expected", 7, "got", len(row))
			continue
		}

		entity := &permissionentities.PermissionEntity{}
		if val, ok := row[0].(int64); ok {
			entity.ID = val
		}
		switch nameVal := row[1].(type) {
		case []byte:
			entity.Name = string(nameVal)
		case string:
			entity.Name = nameVal
		}
		switch resourceVal := row[2].(type) {
		case []byte:
			entity.Resource = string(resourceVal)
		case string:
			entity.Resource = resourceVal
		}
		switch actionVal := row[3].(type) {
		case []byte:
			entity.Action = string(actionVal)
		case string:
			entity.Action = actionVal
		}
		if row[4] != nil {
			switch desc := row[4].(type) {
			case []byte:
				entity.Description = string(desc)
			case string:
				entity.Description = desc
			}
		}
		if row[5] != nil {
			switch cond := row[5].(type) {
			case []byte:
				condStr := string(cond)
				entity.Conditions = &condStr
			case string:
				condStr := cond
				entity.Conditions = &condStr
			}
		}
		switch activeVal := row[6].(type) {
		case int64:
			entity.IsActive = activeVal == 1
		}

		permission, convertErr := permissionconverters.PermissionEntityToDomain(entity)
		if convertErr != nil {
			utils.SetSpanError(ctx, convertErr)
			logger.Error("mysql.permission.list_permissions.convert_error", "row_index", idx, "error", convertErr)
			return permissionrepository.PermissionListResult{}, fmt.Errorf("convert permission entity to domain: %w", convertErr)
		}
		if permission != nil {
			result.Permissions = append(result.Permissions, permission)
		}
	}

	countQuery := `SELECT COUNT(*) FROM permissions p ` + whereClause
	countRows, countErr := pa.Read(ctx, tx, countQuery, args...)
	if countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.permission.list_permissions.count_error", "error", countErr)
		return permissionrepository.PermissionListResult{}, fmt.Errorf("list permissions count: %w", countErr)
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
