package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserRole updates an existing user_role record with new field values
//
// This function updates all fields of a user_role record identified by its ID.
// Used when modifying role assignment, status, activation state, or expiration.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - userRole: UserRoleInterface with ID and fields to update
//
// Returns:
//   - error: sql.ErrNoRows if user_role not found, database errors
//
// Business Rules:
//   - ID must be set and > 0 (identifies user_role to update)
//   - User role must exist in user_roles table
//   - Updates ALL fields: user_id, role_id, is_active, status, expires_at
//   - Does NOT validate foreign key references (DB constraints enforce)
//
// Database Schema:
//   - Table: user_roles
//   - Primary Key: id
//   - Foreign Keys: user_id -> users.id, role_id -> roles.id
//   - Columns: id, user_id, role_id, is_active, status, expires_at, blocked_until
//
// Edge Cases:
//   - Returns sql.ErrNoRows if user_role doesn't exist (service maps to 404)
//   - Invalid user_id/role_id triggers foreign key constraint error
//   - Multiple active roles for same user should be prevented by service layer
//
// Performance:
//   - Single-row UPDATE using PRIMARY KEY (very fast)
//   - Foreign key checks add minimal overhead
//
// Important Notes:
//   - Does NOT deactivate other roles (service layer responsibility)
//   - Does NOT update blocked_until (use BlockUserTemporarily for that)
//   - Transaction managed by service layer
//
// Example:
//
//	userRole := usermodel.NewUserRole()
//	userRole.SetID(userRoleID)
//	userRole.SetStatus(globalmodel.StatusActive)
//	userRole.SetIsActive(true)
//
//	err := adapter.UpdateUserRole(ctx, tx, userRole)
//	if err == sql.ErrNoRows {
//	    // Handle user role not found (404)
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) UpdateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Convert domain model to database entity
	entity, convertErr := userconverters.UserRoleDomainToEntity(userRole)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.user.update_user_role.convert_error", "error", convertErr)
		return fmt.Errorf("convert user role domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.user.update_user_role.empty_entity")
		return nil
	}

	logger = logger.With(
		"user_role_id", entity.ID,
		"user_id", entity.UserID,
		"role_id", entity.RoleID,
	)

	// Update all user_role fields by primary key
	// Note: Foreign key constraints ensure user_id and role_id are valid
	query := `
		UPDATE user_roles 
		SET user_id = ?, role_id = ?, is_active = ?, status = ?, expires_at = ?
		WHERE id = ?
	`

	// Execute update using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		entity.UserID,
		entity.RoleID,
		entity.IsActive,
		entity.Status,
		entity.ExpiresAt,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_role.exec_error", "error", execErr)
		return fmt.Errorf("update user role: %w", execErr)
	}

	// Check if user role exists and was updated
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user role update rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if user role not found (service layer maps to 404)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Debug("mysql.user.update_user_role.success", "rows_affected", rowsAffected)
	return nil
}
