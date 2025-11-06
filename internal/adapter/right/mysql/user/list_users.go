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

// ListUsersWithFilters retrieves a paginated list of users with advanced filtering and role information
//
// This function supports the admin panel user listing with comprehensive filters,
// pagination, and role information via LEFT JOIN. Designed for administrative queries only.
//
// Query Structure:
//   - Base: users table with LEFT JOIN to user_roles and roles
//   - LEFT JOIN ensures users without roles are included (HasRole=false)
//   - Filters only active roles (ur.is_active = 1) to show current role assignment
//   - Returns both user data AND role information in single query (reduces N+1)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for read-only queries)
//   - filter: ListUsersFilter with pagination and optional filters
//
// Returns:
//   - result: ListUsersResult containing users slice and total count
//   - error: Query execution errors, scan errors
//
// Pagination:
//   - Default page=1, limit=20 if not provided
//   - OFFSET calculated as (page-1)*limit
//   - Total count query runs separately (COUNT DISTINCT to handle JOIN duplicates)
//
// Filters Applied (when non-empty):
//   - RoleName/RoleSlug: Filter by role (LIKE for partial match)
//   - RoleStatus: Filter by user_role status (0=pending, 1=active, etc.)
//   - IsSystemRole: Filter system vs custom roles
//   - FullName/CPF/Email/Phone: User field filters (LIKE for partial match)
//   - Deleted: Include/exclude soft-deleted users
//   - ID/BornAt/LastActivity ranges: Date/ID range filters
//
// Performance Considerations:
//   - Uses indexes: user_roles.user_id, user_roles.role_id, roles.slug
//   - COUNT DISTINCT necessary to avoid duplicate counting in JOIN
//   - Consider adding EXPLAIN ANALYZE in dev for query optimization
//
// Edge Cases:
//   - User without role: HasRole=false, role fields are empty
//   - Multiple roles: Only ACTIVE role returned (ur.is_active = 1 filter)
//   - No results: Returns empty slice + total=0 (NOT sql.ErrNoRows)
//
// Security:
//   - Only accessible via admin handlers (permission check in service layer)
//   - Sensitive fields (password hash) returned but masked by DTOs
//
// Example:
//
//	filter := userrepository.ListUsersFilter{
//	    Page: 1,
//	    Limit: 50,
//	    RoleSlug: "owner",
//	    RoleStatus: &activeStatus,
//	    Deleted: &falseValue,
//	}
//	result, err := adapter.ListUsersWithFilters(ctx, nil, filter)
//	// result.Users contains up to 50 owners with active status
//	// result.Total shows total matching users (for pagination UI)
func (ua *UserAdapter) ListUsersWithFilters(ctx context.Context, tx *sql.Tx, filter userrepository.ListUsersFilter) (result userrepository.ListUsersResult, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return result, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Apply default pagination values if not provided
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	var conditions []string
	var args []any

	// Base condition (always true, simplifies conditional logic)
	conditions = append(conditions, "1=1")

	// Build WHERE clause dynamically based on provided filters
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

	// SELECT query with LEFT JOIN to include users without roles
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

	// Combine SELECT with WHERE and ORDER BY (newest first) + pagination
	listQuery := baseSelect + " " + whereClause + " ORDER BY u.id DESC LIMIT ? OFFSET ?"

	// COUNT query to get total matching records (for pagination metadata)
	// Uses COUNT DISTINCT to handle LEFT JOIN duplicates correctly
	countQuery := `SELECT COUNT(DISTINCT u.id)
        FROM users u
        LEFT JOIN user_roles ur ON ur.user_id = u.id AND ur.is_active = 1
        LEFT JOIN roles r ON r.id = ur.role_id
    ` + " " + whereClause

	// Prepare args for list query (filters + limit + offset)
	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	// Execute list query using instrumented adapter (auto-generates metrics + tracing)
	rows, queryErr := ua.QueryContext(ctx, tx, "select", listQuery, listArgs...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.list_users.query_error", "error", queryErr)
		return result, fmt.Errorf("list users query: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entities with role information
	userEntities, scanErr := scanUserEntitiesWithRoles(rows)
	if scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.list_users.scan_error", "error", scanErr)
		return result, fmt.Errorf("scan user entities: %w", scanErr)
	}

	// Pre-allocate slice if we have results (optimization)
	if len(userEntities) > 0 {
		result.Users = make([]usermodel.UserInterface, 0, len(userEntities))
	}

	// Convert entities to domain models
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

	// Execute count query to get total matching records
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
