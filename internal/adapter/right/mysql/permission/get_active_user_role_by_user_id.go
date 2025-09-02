package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetActiveUserRoleByUserID retorna o role ativo único do usuário
func (pa *PermissionAdapter) GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (permissionmodel.UserRoleInterface, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.is_active, ur.status, ur.expires_at
		FROM user_roles ur
		WHERE ur.user_id = ? AND ur.is_active = 1
		LIMIT 1
	`
	// Typed scan in a single round-trip to avoid fragile double-query logic
	var (
		id          int64
		uid         int64
		roleID      int64
		isActiveInt int64
		status      int64
		expiresAt   sql.NullTime
	)

	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&id, &uid, &roleID, &isActiveInt, &status, &expiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Nenhum role ativo encontrado
		}
		slog.Error("mysqlpermissionadapter/GetActiveUserRoleByUserID: error scanning row", "error", err)
		return nil, err
	}

	entity := &permissionentities.UserRoleEntity{
		ID:       id,
		UserID:   uid,
		RoleID:   roleID,
		IsActive: isActiveInt == 1,
		Status:   status,
	}
	if expiresAt.Valid {
		t := expiresAt.Time
		entity.ExpiresAt = &t
	}

	return permissionconverters.UserRoleEntityToDomain(entity), nil
}
