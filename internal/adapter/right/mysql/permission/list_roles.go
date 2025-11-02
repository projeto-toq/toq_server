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

// ListRoles retorna roles com filtros e paginação
func (pa *PermissionAdapter) ListRoles(ctx context.Context, tx *sql.Tx, filter permissionrepository.RoleListFilter) (permissionrepository.RoleListResult, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return permissionrepository.RoleListResult{}, err
	}
	defer spanEnd()

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	var conditions []string
	var args []any
	conditions = append(conditions, "1=1")

	if filter.Name != "" {
		conditions = append(conditions, "r.name LIKE ?")
		args = append(args, filter.Name)
	}
	if filter.Slug != "" {
		conditions = append(conditions, "r.slug LIKE ?")
		args = append(args, filter.Slug)
	}
	if filter.Description != "" {
		conditions = append(conditions, "r.description LIKE ?")
		args = append(args, filter.Description)
	}
	if filter.IsSystemRole != nil {
		conditions = append(conditions, "r.is_system_role = ?")
		if *filter.IsSystemRole {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if filter.IsActive != nil {
		conditions = append(conditions, "r.is_active = ?")
		if *filter.IsActive {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if filter.IDFrom != nil {
		conditions = append(conditions, "r.id >= ?")
		args = append(args, *filter.IDFrom)
	}
	if filter.IDTo != nil {
		conditions = append(conditions, "r.id <= ?")
		args = append(args, *filter.IDTo)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	selectQuery := `SELECT id, name, slug, description, is_system_role, is_active
        FROM roles r
    ` + " " + whereClause + " ORDER BY name ASC LIMIT ? OFFSET ?"

	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	rows, readErr := pa.QueryContext(ctx, tx, "select", selectQuery, listArgs...)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.permission.list_roles.read_error", "error", readErr)
		return permissionrepository.RoleListResult{}, fmt.Errorf("list roles read: %w", readErr)
	}
	defer rows.Close()

	rowEntities, rowsErr := rowsToEntities(rows)
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.list_roles.rows_to_entities_error", "error", rowsErr)
		return permissionrepository.RoleListResult{}, fmt.Errorf("list roles rows to entities: %w", rowsErr)
	}

	result := permissionrepository.RoleListResult{}
	if len(rowEntities) > 0 {
		result.Roles = make([]permissionmodel.RoleInterface, 0, len(rowEntities))
	}

	for _, row := range rowEntities {
		if len(row) != 6 {
			logger.Warn("mysql.permission.list_roles.columns_mismatch", "expected", 6, "got", len(row))
			continue
		}
		entity := &permissionentities.RoleEntity{}
		if val, ok := row[0].(int64); ok {
			entity.ID = val
		}
		switch nameVal := row[1].(type) {
		case []byte:
			entity.Name = string(nameVal)
		case string:
			entity.Name = nameVal
		}
		switch slugVal := row[2].(type) {
		case []byte:
			entity.Slug = string(slugVal)
		case string:
			entity.Slug = slugVal
		}
		if row[3] != nil {
			switch desc := row[3].(type) {
			case []byte:
				entity.Description = string(desc)
			case string:
				entity.Description = desc
			}
		}
		switch sysVal := row[4].(type) {
		case int64:
			entity.IsSystemRole = sysVal == 1
		}
		switch activeVal := row[5].(type) {
		case int64:
			entity.IsActive = activeVal == 1
		}

		role := permissionconverters.RoleEntityToDomain(entity)
		if role != nil {
			result.Roles = append(result.Roles, role)
		}
	}

	countQuery := `SELECT COUNT(*) FROM roles r ` + whereClause
	countRow := pa.QueryRowContext(ctx, tx, "select", countQuery, args...)
	var total int64
	if countErr := countRow.Scan(&total); countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.permission.list_roles.count_error", "error", countErr)
		return permissionrepository.RoleListResult{}, fmt.Errorf("list roles count: %w", countErr)
	}
	result.Total = total

	return result, nil
}
