package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetUserRoleByUserIDAndRoleID busca um user_role específico pela combinação user_id + role_id
func (pa *PermissionAdapter) GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (permissionmodel.UserRoleInterface, error) {
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

	err := tx.QueryRowContext(ctx, query, userID, roleID).Scan(
		&id, &uid, &roleIDOut, &isActiveInt, &status, &expiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("mysqlpermissionadapter/GetUserRoleByUserIDAndRoleID: error scanning row", "error", err)
		return nil, err
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

	return permissionconverters.UserRoleEntityToDomain(entity), nil
}
