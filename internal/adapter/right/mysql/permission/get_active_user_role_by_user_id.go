package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetActiveUserRoleByUserID retorna o role ativo único do usuário
func (pa *PermissionAdapter) GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (permissionmodel.UserRoleInterface, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.is_active, ur.expires_at
		FROM user_roles ur
		WHERE ur.user_id = ? AND ur.is_active = 1
		LIMIT 1
	`

	row, err := pa.ReadRow(ctx, tx, query, userID)
	if err != nil {
		return nil, err
	}

	if row == nil {
		return nil, nil // Nenhum role ativo encontrado
	}

	if len(row) != 5 {
		return nil, fmt.Errorf("unexpected number of columns: expected 5, got %d", len(row))
	}

	entity := &permissionentities.UserRoleEntity{
		ID:       row[0].(int64),
		UserID:   row[1].(int64),
		RoleID:   row[2].(int64),
		IsActive: row[3].(int64) == 1,
	}

	// Handle expires_at (pode ser NULL)
	if row[4] != nil {
		if expiresAt, ok := row[4].(time.Time); ok {
			entity.ExpiresAt = &expiresAt
		}
	}

	return permissionconverters.UserRoleEntityToDomain(entity), nil
}
