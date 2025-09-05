package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetActiveUserRoleByUserID retorna o role ativo único do usuário
func (pa *PermissionAdapter) GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (permissionmodel.UserRoleInterface, error) {
	// tracing mínimo no adapter, conforme guia
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	// Query com JOIN em roles para popular o Role no UserRole
	query := `
		SELECT 
			ur.id,
			ur.user_id,
			ur.role_id,
			ur.is_active,
			ur.status,
			ur.expires_at,
			r.id,
			r.slug,
			r.name,
			r.description,
			r.is_system_role,
			r.is_active
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ?
		  AND ur.is_active = 1
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		  AND r.is_active = 1
		ORDER BY ur.id DESC
		LIMIT 1`

	// Typed scan em uma única consulta
	var (
		id          int64
		uid         int64
		roleID      int64
		isActiveInt int64
		status      int64
		expiresAt   sql.NullTime

		rID          int64
		rSlug        string
		rName        string
		rDescription sql.NullString
		rIsSystemInt int64
		rIsActiveInt int64
	)

	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&id,
		&uid,
		&roleID,
		&isActiveInt,
		&status,
		&expiresAt,
		&rID,
		&rSlug,
		&rName,
		&rDescription,
		&rIsSystemInt,
		&rIsActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Nenhum role ativo encontrado
			return nil, nil
		}
		slog.Error("mysqlpermissionadapter/GetActiveUserRoleByUserID: error scanning row", "error", err)
		return nil, err
	}

	// Monta entidades tipadas
	userRoleEntity := &permissionentities.UserRoleEntity{
		ID:       id,
		UserID:   uid,
		RoleID:   roleID,
		IsActive: isActiveInt == 1,
		Status:   status,
	}
	if expiresAt.Valid {
		t := expiresAt.Time
		userRoleEntity.ExpiresAt = &t
	}

	roleEntity := &permissionentities.RoleEntity{
		ID:   rID,
		Name: rName,
		Slug: rSlug,
		Description: func() string {
			if rDescription.Valid {
				return rDescription.String
			}
			return ""
		}(),
		IsSystemRole: rIsSystemInt == 1,
		IsActive:     rIsActiveInt == 1,
	}

	// Converte para domínio e associa Role ao UserRole
	userRole := permissionconverters.UserRoleEntityToDomain(userRoleEntity)
	if userRole != nil {
		role := permissionconverters.RoleEntityToDomain(roleEntity)
		userRole.SetRole(role)
	}

	return userRole, nil
}
