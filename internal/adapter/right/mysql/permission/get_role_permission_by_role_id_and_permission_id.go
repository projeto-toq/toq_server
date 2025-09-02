package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetRolePermissionByRoleIDAndPermissionID busca um role_permission específico pela combinação role_id + permission_id
func (pa *PermissionAdapter) GetRolePermissionByRoleIDAndPermissionID(ctx context.Context, tx *sql.Tx, roleID, permissionID int64) (permissionmodel.RolePermissionInterface, error) {
	query := `
		SELECT id, role_id, permission_id, granted, conditions
		FROM role_permissions 
		WHERE role_id = ? AND permission_id = ?
		LIMIT 1
	`

	var (
		id              int64
		roleIDOut       int64
		permissionIDOut int64
		grantedInt      int64
		conditions      sql.NullString
	)

	err := tx.QueryRowContext(ctx, query, roleID, permissionID).Scan(
		&id, &roleIDOut, &permissionIDOut, &grantedInt, &conditions,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("mysqlpermissionadapter/GetRolePermissionByRoleIDAndPermissionID: error scanning row", "error", err)
		return nil, err
	}

	entity := &permissionentities.RolePermissionEntity{
		ID:           id,
		RoleID:       roleIDOut,
		PermissionID: permissionIDOut,
		Granted:      grantedInt == 1,
	}
	if conditions.Valid {
		v := conditions.String
		entity.Conditions = &v
	}

	return permissionconverters.RolePermissionEntityToDomain(entity), nil
}
