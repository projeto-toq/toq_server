package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserByNationalID retrieves a user with their active role by national ID (CPF or CNPJ).
//
// This function performs the same optimized JOIN as GetUserByID but filters by national_id instead.
// It does NOT filter by deleted status, allowing authentication checks to detect deleted accounts
// (prevents account enumeration attacks).
//
// Query Strategy:
//   - LEFT JOIN users + user_roles + roles
//   - No deleted filter (returns deleted users for security checks)
//   - Filters user_roles by is_active = 1
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - nationalID: User's CPF (11 digits) or CNPJ (14 digits) without formatting
//
// Returns:
//   - user: UserInterface with all fields including ActiveRole if exists
//   - error: sql.ErrNoRows if not found, or database errors
//
// Business Rules:
//   - National ID is UNIQUE constraint in database
//   - Does NOT filter by deleted (returns even deleted users)
//   - Service layer decides if deleted users are acceptable for the operation
//
// Security Considerations:
//   - Used for authentication (SignIn, password reset)
//   - Returns deleted users to prevent account enumeration
//   - Service logs authentication attempts for deleted accounts
//
// Performance:
//   - Single query replaces 2-query pattern
//   - Critical path optimization (SignIn uses this function)
func (ua *UserAdapter) GetUserByNationalID(ctx context.Context, tx *sql.Tx, nationalID string) (user usermodel.UserInterface, err error) {
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
	// Note: NO deleted filter (security requirement for authentication)
	query := `
		SELECT 
			u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state,
			u.creci_validity, u.born_at, u.phone_number, u.email, u.zip_code, u.street, 
			u.number, u.complement, u.neighborhood, u.city, u.state, u.password, 
			u.opt_status, u.last_activity_at, u.deleted,
			ur.id AS user_role_id, ur.user_id AS user_role_user_id, ur.role_id AS user_role_role_id,
			ur.is_active AS user_role_is_active, ur.status AS user_role_status, 
			ur.expires_at AS user_role_expires_at, ur.blocked_until AS user_role_blocked_until,
			r.id AS role_id, r.slug AS role_slug, r.name AS role_name, 
			r.description AS role_description, r.is_system_role AS role_is_system_role, 
			r.is_active AS role_is_active
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id AND ur.is_active = 1
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.national_id = ?
	`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, nationalID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_by_national_id.query_error", "error", queryErr)
		return nil, fmt.Errorf("query user by national_id with role: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entity
	entities, err := scanUserWithRoleEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_national_id.scan_error", "error", err)
		return nil, fmt.Errorf("scan user with role rows: %w", err)
	}

	// Handle no results
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: unique constraint should prevent multiple rows
	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple active roles found for national_id: %s", nationalID)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_national_id.multiple_roles_error",
			"national_id", nationalID, "count", len(entities), "error", errMultiple)
		return nil, errMultiple
	}

	// Convert entity to domain model
	user, convErr := userconverters.UserWithRoleEntityToDomain(entities[0])
	if convErr != nil {
		utils.SetSpanError(ctx, convErr)
		logger.Error("mysql.user.get_user_by_national_id.conversion_error", "error", convErr)
		return nil, fmt.Errorf("convert user with role to domain: %w", convErr)
	}

	return user, nil
}
