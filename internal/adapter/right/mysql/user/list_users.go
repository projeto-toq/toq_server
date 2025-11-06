package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	userrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListUsersWithFilters retorna usuários com filtros e paginação para o painel admin.
func (ua *UserAdapter) ListUsersWithFilters(ctx context.Context, tx *sql.Tx, filter userrepository.ListUsersFilter) (result userrepository.ListUsersResult, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return result, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	var conditions []string
	var args []any

	conditions = append(conditions, "1=1")

	if filter.RoleName != "" {
		conditions = append(conditions, "r.name LIKE ?")
		args = append(args, filter.RoleName)
	}
	if filter.RoleSlug != "" {
		conditions = append(conditions, "r.slug LIKE ?")
		args = append(args, filter.RoleSlug)
	}
	if filter.RoleStatus != nil {
		conditions = append(conditions, "ur.status = ?")
		args = append(args, int(*filter.RoleStatus))
	}
	if filter.IsSystemRole != nil {
		conditions = append(conditions, "r.is_system_role = ?")
		if *filter.IsSystemRole {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if filter.FullName != "" {
		conditions = append(conditions, "u.full_name LIKE ?")
		args = append(args, filter.FullName)
	}
	if filter.CPF != "" {
		conditions = append(conditions, "u.national_id LIKE ?")
		args = append(args, filter.CPF)
	}
	if filter.Email != "" {
		conditions = append(conditions, "u.email LIKE ?")
		args = append(args, filter.Email)
	}
	if filter.PhoneNumber != "" {
		conditions = append(conditions, "u.phone_number LIKE ?")
		args = append(args, filter.PhoneNumber)
	}
	if filter.Deleted != nil {
		conditions = append(conditions, "u.deleted = ?")
		if *filter.Deleted {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if filter.IDFrom != nil {
		conditions = append(conditions, "u.id >= ?")
		args = append(args, *filter.IDFrom)
	}
	if filter.IDTo != nil {
		conditions = append(conditions, "u.id <= ?")
		args = append(args, *filter.IDTo)
	}
	if filter.BornAtFrom != nil {
		conditions = append(conditions, "u.born_at >= ?")
		args = append(args, *filter.BornAtFrom)
	}
	if filter.BornAtTo != nil {
		conditions = append(conditions, "u.born_at <= ?")
		args = append(args, *filter.BornAtTo)
	}
	if filter.LastActivityFrom != nil {
		conditions = append(conditions, "u.last_activity_at >= ?")
		args = append(args, *filter.LastActivityFrom)
	}
	if filter.LastActivityTo != nil {
		conditions = append(conditions, "u.last_activity_at <= ?")
		args = append(args, *filter.LastActivityTo)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	baseSelect := `SELECT 
        u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state, u.creci_validity,
        u.born_at, u.phone_number, u.email, u.zip_code, u.street, u.number, u.complement,
        u.neighborhood, u.city, u.state, u.password, u.opt_status, u.last_activity_at, u.deleted, u.last_signin_attempt,
        ur.id AS active_user_role_id, ur.role_id, ur.status, ur.is_active,
        r.id AS role_id, r.name, r.slug, r.description, r.is_system_role, r.is_active
        FROM users u
        LEFT JOIN user_roles ur ON ur.user_id = u.id AND ur.is_active = 1
        LEFT JOIN roles r ON r.id = ur.role_id
    `

	listQuery := baseSelect + " " + whereClause + " ORDER BY u.id DESC LIMIT ? OFFSET ?"

	countQuery := `SELECT COUNT(DISTINCT u.id)
        FROM users u
        LEFT JOIN user_roles ur ON ur.user_id = u.id AND ur.is_active = 1
        LEFT JOIN roles r ON r.id = ur.role_id
    ` + " " + whereClause

	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	rows, queryErr := ua.QueryContext(ctx, tx, "select", listQuery, listArgs...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.list_users.query_error", "error", queryErr)
		return result, fmt.Errorf("list users query: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entities (type-safe scanning)
	userEntities, scanErr := scanUserEntitiesWithRoles(rows)
	if scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.list_users.scan_error", "error", scanErr)
		return result, fmt.Errorf("scan user entities: %w", scanErr)
	}

	if len(userEntities) > 0 {
		result.Users = make([]usermodel.UserInterface, 0, len(userEntities))
	}

	for _, entity := range userEntities {
		// Convert user entity to domain model
		user := userconverters.UserEntityToDomain(entity.User)

		// Attach role information if present (LEFT JOIN may return NULL)
		if entity.HasRole {
			userRole := permissionmodel.NewUserRole()
			userRole.SetID(entity.UserRoleID)
			userRole.SetUserID(user.GetID())
			userRole.SetRoleID(entity.RoleID)
			userRole.SetStatus(permissionmodel.UserRoleStatus(entity.RoleStatus))
			userRole.SetIsActive(entity.RoleIsActive)

			// Attach role details
			role := permissionmodel.NewRole()
			role.SetID(entity.RoleID)
			role.SetName(entity.RoleName)
			role.SetSlug(entity.RoleSlug)
			role.SetDescription(entity.RoleDescription)
			role.SetIsSystemRole(entity.RoleIsSystemRole)
			role.SetIsActive(entity.RoleActive)

			userRole.SetRole(role)
			user.SetActiveRole(userRole)
		}

		result.Users = append(result.Users, user)
	}

	countArgs := append([]any{}, args...)
	row := ua.QueryRowContext(ctx, tx, "select", countQuery, countArgs...)

	var total int64
	if scanErr := row.Scan(&total); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.list_admin.count_scan_error", "error", scanErr)
		return result, fmt.Errorf("list admin users count scan: %w", scanErr)
	}
	result.Total = total

	return result, nil
}
