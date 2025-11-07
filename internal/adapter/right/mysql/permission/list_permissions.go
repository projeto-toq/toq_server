package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListPermissions retorna permissões com filtros e paginação
func (p *PermissionAdapter) ListPermissions(ctx context.Context, tx *sql.Tx, filter permissionrepository.PermissionListFilter) (permissionrepository.PermissionListResult, error) {
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

	query := `SELECT id, name, action, description, is_active
		FROM permissions p ` + " " + whereClause + " ORDER BY action ASC LIMIT ? OFFSET ?"

	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	rows, readErr := p.QueryContext(ctx, tx, "select", query, listArgs...)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.permission.list_permissions.read_error", "error", readErr)
		return permissionrepository.PermissionListResult{}, fmt.Errorf("list permissions read: %w", readErr)
	}
	defer rows.Close()

	rowEntities, rowsErr := rowsToEntities(rows)
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.list_permissions.rows_to_entities_error", "error", rowsErr)
		return permissionrepository.PermissionListResult{}, fmt.Errorf("list permissions rows to entities: %w", rowsErr)
	}

	result := permissionrepository.PermissionListResult{}
	if len(rowEntities) > 0 {
		result.Permissions = make([]permissionmodel.PermissionInterface, 0, len(rowEntities))
	}

	for idx, row := range rowEntities {
		if len(row) != 5 {
			logger.Warn("mysql.permission.list_permissions.columns_mismatch", "row_index", idx, "expected", 5, "got", len(row))
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
		switch actionVal := row[2].(type) {
		case []byte:
			entity.Action = string(actionVal)
		case string:
			entity.Action = actionVal
		}
		if row[3] != nil {
			switch desc := row[3].(type) {
			case []byte:
				d := string(desc)
				entity.Description = &d
			case string:
				d := desc
				entity.Description = &d
			}
		}
		switch activeVal := row[4].(type) {
		case int64:
			entity.IsActive = activeVal == 1
		case bool:
			entity.IsActive = activeVal
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
	countRow := p.QueryRowContext(ctx, tx, "select", countQuery, args...)
	var total int64
	if countErr := countRow.Scan(&total); countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.permission.list_permissions.count_error", "error", countErr)
		return permissionrepository.PermissionListResult{}, fmt.Errorf("list permissions count: %w", countErr)
	}
	result.Total = total

	return result, nil
}
