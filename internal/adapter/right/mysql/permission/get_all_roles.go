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

// GetAllRoles busca todos os roles
func (pa *PermissionAdapter) GetAllRoles(ctx context.Context, tx *sql.Tx) (roles []permissionmodel.RoleInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	query := `
		SELECT id, name, slug, description, is_system_role, is_active
		FROM roles 
		ORDER BY name
	`

	results, err := pa.Read(ctx, tx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_all_roles.read_error", "error", err)
		return nil, fmt.Errorf("get all roles read: %w", err)
	}

	roles = make([]permissionmodel.RoleInterface, 0, len(results))
	for index, row := range results {
		if len(row) != 6 {
			errColumns := fmt.Errorf("unexpected number of columns: expected 6, got %d", len(row))
			utils.SetSpanError(ctx, errColumns)
			logger.Error("mysql.permission.get_all_roles.columns_mismatch", "row_index", index, "error", errColumns)
			return nil, errColumns
		}

		entity := &permissionentities.RoleEntity{
			ID:           row[0].(int64),
			Name:         string(row[1].([]byte)),
			Slug:         string(row[2].([]byte)),
			IsSystemRole: row[4].(int64) == 1,
			IsActive:     row[5].(int64) == 1,
		}

		// description (pode ser NULL)
		if row[3] != nil {
			entity.Description = string(row[3].([]byte))
		}

		role := permissionconverters.RoleEntityToDomain(entity)
		if role != nil {
			roles = append(roles, role)
		}
	}

	logger.Debug("mysql.permission.get_all_roles.success", "count", len(roles))
	return roles, nil
}
