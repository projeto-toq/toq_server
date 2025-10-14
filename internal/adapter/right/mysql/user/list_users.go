package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
		args = append(args, "%"+filter.RoleName+"%")
	}
	if filter.RoleSlug != "" {
		conditions = append(conditions, "r.slug = ?")
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
		args = append(args, "%"+filter.FullName+"%")
	}
	if filter.CPF != "" {
		conditions = append(conditions, "u.national_id = ?")
		args = append(args, filter.CPF)
	}
	if filter.Email != "" {
		conditions = append(conditions, "u.email LIKE ?")
		args = append(args, "%"+filter.Email+"%")
	}
	if filter.PhoneNumber != "" {
		conditions = append(conditions, "u.phone_number LIKE ?")
		args = append(args, "%"+filter.PhoneNumber+"%")
	}
	if filter.Deleted != nil {
		conditions = append(conditions, "u.deleted = ?")
		if *filter.Deleted {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
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

	entities, readErr := ua.Read(ctx, tx, listQuery, listArgs...)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.user.list_admin.read_error", "error", readErr)
		return result, fmt.Errorf("list admin users read: %w", readErr)
	}

	if len(entities) > 0 {
		result.Users = make([]usermodel.UserInterface, 0, len(entities))
	}

	for _, row := range entities {
		if len(row) < 32 {
			logger.Warn("mysql.user.list_admin.unexpected_columns", "expected", 32, "got", len(row))
			continue
		}

		userEntitySlice := row[:22]
		user, convertErr := userconverters.UserEntityToDomain(userEntitySlice)
		if convertErr != nil {
			utils.SetSpanError(ctx, convertErr)
			logger.Error("mysql.user.list_admin.convert_error", "error", convertErr)
			return result, fmt.Errorf("convert user entity: %w", convertErr)
		}

		if row[22] != nil {
			userRole := permissionmodel.NewUserRole()
			switch id := row[22].(type) {
			case int64:
				userRole.SetID(id)
			}
			userRole.SetUserID(user.GetID())
			switch roleID := row[23].(type) {
			case int64:
				userRole.SetRoleID(roleID)
			}
			switch statusVal := row[24].(type) {
			case int64:
				userRole.SetStatus(permissionmodel.UserRoleStatus(statusVal))
			}
			switch activeVal := row[25].(type) {
			case int64:
				userRole.SetIsActive(activeVal == 1)
			}

			if row[26] != nil {
				role := permissionmodel.NewRole()
				switch val := row[26].(type) {
				case int64:
					role.SetID(val)
				}
				if row[27] != nil {
					switch nameVal := row[27].(type) {
					case []byte:
						role.SetName(string(nameVal))
					case string:
						role.SetName(nameVal)
					}
				}
				if row[28] != nil {
					switch slugVal := row[28].(type) {
					case []byte:
						role.SetSlug(string(slugVal))
					case string:
						role.SetSlug(slugVal)
					}
				}
				if row[29] != nil {
					switch desc := row[29].(type) {
					case []byte:
						role.SetDescription(string(desc))
					case string:
						role.SetDescription(desc)
					}
				}
				if row[30] != nil {
					switch flag := row[30].(type) {
					case int64:
						role.SetIsSystemRole(flag == 1)
					}
				}
				if row[31] != nil {
					switch flag := row[31].(type) {
					case int64:
						role.SetIsActive(flag == 1)
					}
				}
				userRole.SetRole(role)
			}
			user.SetActiveRole(userRole)
		}

		result.Users = append(result.Users, user)
	}

	countArgs := append([]any{}, args...)
	countRows, countErr := ua.Read(ctx, tx, countQuery, countArgs...)
	if countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.user.list_admin.count_error", "error", countErr)
		return result, fmt.Errorf("list admin users count: %w", countErr)
	}
	if len(countRows) > 0 && len(countRows[0]) > 0 {
		switch total := countRows[0][0].(type) {
		case int64:
			result.Total = total
		case []byte:
			if val, convErr := strconv.ParseInt(string(total), 10, 64); convErr == nil {
				result.Total = val
			}
		}
	}

	return result, nil
}
