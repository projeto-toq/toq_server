package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserRoleByUserIDAndRoleID busca um user_role específico pela combinação user_id + role_id
func (ua *UserAdapter) GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

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

	row := ua.QueryRowContext(ctx, tx, "select", query, userID, roleID)
	err = row.Scan(
		&id, &uid, &roleIDOut, &isActiveInt, &status, &expiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_user_role_by_user_id_and_role_id.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_user_role_by_user_id_and_role_id.scan_error", "error", err)
		return nil, fmt.Errorf("get user role by user id and role id scan: %w", err)
	}

	entity := &userentity.UserRoleEntity{
		ID:       id,
		UserID:   userID,
		RoleID:   roleIDOut,
		IsActive: isActiveInt == 1,
		Status:   status,
	}
	if expiresAt.Valid {
		entity.ExpiresAt = sql.NullTime{
			Time:  expiresAt.Time,
			Valid: true,
		}
	}

	userRole, convertErr := userconverters.UserRoleEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.user.get_user_role_by_user_id_and_role_id.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert user role entity to domain: %w", convertErr)
	}

	logger.Debug("mysql.permission.get_user_role_by_user_id_and_role_id.success")
	return userRole, nil
}
