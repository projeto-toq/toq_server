package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetActiveUserRoleByUserID retorna o role ativo único do usuário
func (ua *UserAdapter) GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

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

	row := ua.QueryRowContext(ctx, tx, "select", query, userID)
	err = row.Scan(
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
			logger.Debug("mysql.permission.get_active_user_role_by_user_id.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_active_user_role_by_user_id.scan_error", "error", err)
		return nil, fmt.Errorf("get active user role by user id scan: %w", err)
	}

	// Monta entidades tipadas
	userRoleEntity := &userentity.UserRoleEntity{
		ID:       id,
		UserID:   userID,
		RoleID:   roleID,
		IsActive: isActiveInt == 1,
		Status:   status,
	}
	if expiresAt.Valid {
		userRoleEntity.ExpiresAt = sql.NullTime{
			Time:  expiresAt.Time,
			Valid: true,
		}
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
	userRole, convertErr := userconverters.UserRoleEntityToDomain(userRoleEntity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_active_user_role_by_user_id.convert_user_role_error", "error", convertErr)
		return nil, fmt.Errorf("convert active user role entity to domain: %w", convertErr)
	}
	if userRole != nil {
		role := permissionconverters.RoleEntityToDomain(roleEntity)
		if role != nil {
			userRole.SetRole(role)
		}
	}

	logger.Debug("mysql.permission.get_active_user_role_by_user_id.success", "role_id", roleID)
	return userRole, nil
}
