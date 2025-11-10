package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateUserRole creates a new user-role association in the database
//
// This function inserts a new record in the user_roles table, assigning a role
// to a user with specific status, activation state, and optional expiration.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - role assignment must be transactional)
//   - userRole: UserRoleInterface with all required fields populated (ID will be set)
//
// Returns:
//   - result: UserRoleInterface with ID populated from auto-generated primary key
//   - error: Database errors (constraint violations, connection issues)
//
// Side Effects:
//   - Modifies userRole object by setting ID to the newly inserted row's primary key
//
// Business Rules:
//   - User and role must exist (foreign key constraints)
//   - Duplicate user-role pairs are prevented by unique constraint (uk_user_roles)
//   - is_active defaults to 1 if not set
//   - status defaults to 0 if not set
//   - expires_at can be NULL for permanent role assignments
//
// Database Constraints:
//   - FK: user_id REFERENCES users(id) ON DELETE CASCADE
//   - FK: role_id REFERENCES roles(id) ON DELETE CASCADE
//   - UNIQUE: (user_id, role_id)
//
// Usage Example:
//
//	userRole := usermodel.NewUserRole()
//	userRole.SetUserID(123)
//	userRole.SetRoleID(456)
//	userRole.SetIsActive(true)
//	userRole.SetStatus(0)
//	created, err := adapter.CreateUserRole(ctx, tx, userRole)
//	if err != nil {
//	    // Handle duplicate or FK violation
//	}
func (ua *UserAdapter) CreateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) (result usermodel.UserRoleInterface, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Convert domain model to database entity
	entity, convertErr := userconverters.UserRoleDomainToEntity(userRole)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.user.create_user_role.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert user role domain to entity: %w", convertErr)
	}

	// Safety check: ensure converter returned valid entity
	if entity == nil {
		logger.Warn("mysql.user.create_user_role.empty_entity")
		return nil, nil
	}

	logger = logger.With(
		"user_id", entity.UserID,
		"role_id", entity.RoleID,
	)

	// Insert new user-role association
	query := `
		INSERT INTO user_roles (user_id, role_id, is_active, status, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`

	// Execute insert using instrumented adapter
	resultExec, execErr := ua.ExecContext(ctx, tx, "insert", query,
		entity.UserID,
		entity.RoleID,
		entity.IsActive,
		entity.Status,
		entity.ExpiresAt,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.create_user_role.exec_error", "error", execErr)
		return nil, fmt.Errorf("create user role: %w", execErr)
	}

	// Retrieve auto-generated primary key
	id, lastErr := resultExec.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.user.create_user_role.last_insert_id_error", "error", lastErr)
		return nil, fmt.Errorf("user role last insert id: %w", lastErr)
	}

	// Update user role with generated ID (side effect)
	userRole.SetID(id)
	logger.Debug("mysql.user.create_user_role.success", "user_role_id", id)

	return userRole, nil
}
