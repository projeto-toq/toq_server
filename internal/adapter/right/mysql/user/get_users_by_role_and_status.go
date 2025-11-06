package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUsersByRoleAndStatus returns users that have an active user_role with the given role slug and status
func (ua *UserAdapter) GetUsersByRoleAndStatus(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleSlug, status permissionmodel.UserRoleStatus) ([]usermodel.UserInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query joins users with roles and user_roles to filter by active role and status
	query := `
        SELECT u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state, u.creci_validity,
               u.born_at, u.phone_number, u.email, u.zip_code, u.street, u.number, u.complement,
               u.neighborhood, u.city, u.state, u.password, u.opt_status, u.last_activity_at, u.deleted, u.last_signin_attempt
          FROM users u
          JOIN user_roles ur ON ur.user_id = u.id AND ur.is_active = 1 AND ur.status = ?
          JOIN roles r ON r.id = ur.role_id AND r.slug = ?
         WHERE u.deleted = 0`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, int(status), role)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_users_by_role.query_error", "error", queryErr, "role", role, "status", status)
		return nil, fmt.Errorf("get users by role and status query: %w", queryErr)
	}
	defer rows.Close()

	entities, qerr := scanUserEntities(rows)
	if qerr != nil {
		utils.SetSpanError(ctx, qerr)
		logger.Error("mysql.user.get_users_by_role.scan_error", "error", qerr)
		return nil, fmt.Errorf("scan users by role rows: %w", qerr)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	users := make([]usermodel.UserInterface, 0, len(entities))
	for _, e := range entities {
		u := userconverters.UserEntityToDomain(e)
		users = append(users, u)
	}
	return users, nil
}
