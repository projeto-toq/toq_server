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

// GetUserPermissions busca todas as permissões efetivas de um usuário através de seus roles
func (p *PermissionAdapter) GetUserPermissions(ctx context.Context, tx *sql.Tx, userID int64) (permissions []permissionmodel.PermissionInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID)

	query := `
		SELECT DISTINCT 
			p.id, p.name, p.action, p.description, p.is_active
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id AND rp.granted = 1
		INNER JOIN roles r ON rp.role_id = r.id AND r.is_active = 1
		INNER JOIN user_roles ur ON r.id = ur.role_id AND ur.is_active = 1
		WHERE ur.user_id = ? 
		  AND p.is_active = 1
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		ORDER BY p.action
	`

	rows, readErr := p.QueryContext(ctx, tx, "select", query, userID)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.permission.get_user_permissions.read_error", "error", readErr)
		return nil, fmt.Errorf("get user permissions read: %w", readErr)
	}
	defer rows.Close()

	rowEntities, rowsErr := rowsToEntities(rows)
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.get_user_permissions.rows_to_entities_error", "error", rowsErr)
		return nil, fmt.Errorf("get user permissions rows to entities: %w", rowsErr)
	}

	permissions = make([]permissionmodel.PermissionInterface, 0, len(rowEntities))
	for index, row := range rowEntities {
		if len(row) != 5 {
			errColumns := fmt.Errorf("unexpected number of columns: expected 5, got %d", len(row))
			utils.SetSpanError(ctx, errColumns)
			logger.Error("mysql.permission.get_user_permissions.columns_mismatch", "row_index", index, "error", errColumns)
			return nil, errColumns
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
			logger.Error("mysql.permission.get_user_permissions.convert_error", "row_index", index, "error", convertErr)
			return nil, fmt.Errorf("convert permission entity to domain: %w", convertErr)
		}
		if permission != nil {
			permissions = append(permissions, permission)
		}
	}

	logger.Debug("mysql.permission.get_user_permissions.success", "count", len(permissions))
	return permissions, nil
}
