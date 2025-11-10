package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserByPhoneNumber retrieves a user with their active role by phone number.
//
// Identical strategy to GetUserByNationalID but filters by phone_number.
// Does NOT filter by deleted status for security reasons.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - phoneNumber: User's phone in E.164 format (e.g., "+5511999999999")
//
// Returns:
//   - user: UserInterface with all fields including ActiveRole
//   - error: sql.ErrNoRows if not found, or database errors
//
// Security Note:
//   - Returns deleted users to prevent account enumeration
//   - Service layer handles deleted account logic
func (ua *UserAdapter) GetUserByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (user usermodel.UserInterface, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Optimized query: JOIN fetches user + active role
	query := `
		SELECT 
			u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state,
			u.creci_validity, u.born_at, u.phone_number, u.email, u.zip_code, u.street, 
			u.number, u.complement, u.neighborhood, u.city, u.state, u.password, 
			u.opt_status, u.last_activity_at, u.deleted, u.last_signin_attempt,
			ur.id AS user_role_id, ur.user_id AS user_role_user_id, ur.role_id AS user_role_role_id,
			ur.is_active AS user_role_is_active, ur.status AS user_role_status, 
			ur.expires_at AS user_role_expires_at, ur.blocked_until AS user_role_blocked_until,
			r.id AS role_id, r.slug AS role_slug, r.name AS role_name, 
			r.description AS role_description, r.is_system_role AS role_is_system_role, 
			r.is_active AS role_is_active
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id AND ur.is_active = 1
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.phone_number = ?
	`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, phoneNumber)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_by_phone.query_error", "error", queryErr)
		return nil, fmt.Errorf("query user by phone with role: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entity
	entities, err := scanUserWithRoleEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_phone.scan_error", "error", err)
		return nil, fmt.Errorf("scan user with role rows: %w", err)
	}

	// Handle no results
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: unique constraint should prevent multiple rows
	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple active roles found for phone: %s", phoneNumber)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_phone.multiple_roles_error",
			"phone", phoneNumber, "count", len(entities), "error", errMultiple)
		return nil, errMultiple
	}

	// Convert entity to domain model
	user, convErr := userconverters.UserWithRoleEntityToDomain(entities[0])
	if convErr != nil {
		utils.SetSpanError(ctx, convErr)
		logger.Error("mysql.user.get_user_by_phone.conversion_error", "error", convErr)
		return nil, fmt.Errorf("convert user with role to domain: %w", convErr)
	}

	return user, nil
}
