package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserByID updates user data for an existing user
//
// This function updates all user fields EXCEPT password and last_activity_at
// which have dedicated update functions for security and performance reasons.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - user: UserInterface with ID and fields to update
//
// Returns:
//   - error: sql.ErrNoRows if user not found, constraint violations, database errors
//
// Business Rules:
//   - ID must be set and > 0 (identifies user to update)
//   - User must exist and not be deleted (WHERE id = ? implicit deleted check)
//   - Email, phone, and national ID uniqueness is enforced (may return constraint violation)
//
// Fields NOT Updated:
//   - password: Use UpdateUserPasswordByID()
//   - last_activity_at: Use UpdateUserLastActivity() or batch update
//   - id: Primary key is immutable
//
// Edge Cases:
//   - Returns sql.ErrNoRows if user deleted or doesn't exist
//   - Constraint violation errors propagated to service layer
func (ua *UserAdapter) UpdateUserByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// SQL update for all fields except password and last_activity_at
	query := `UPDATE users SET
			full_name = ?, nick_name = ?, national_id = ?, creci_number = ?, creci_state = ?, creci_validity = ?,
			born_at = ?, phone_number = ?, email = ?, zip_code = ?, street = ?, number = ?, complement = ?, neighborhood = ?, 
			city = ?, state = ?, opt_status = ?, deleted = ?
			WHERE id = ?`

	// Convert domain model to database entity
	entity := userconverters.UserDomainToEntity(user)

	// Execute update using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		entity.FullName,
		entity.NickName,
		entity.NationalID,
		entity.CreciNumber,
		entity.CreciState,
		entity.CreciValidity,
		entity.BornAt,
		entity.PhoneNumber,
		entity.Email,
		entity.ZipCode,
		entity.Street,
		entity.Number,
		entity.Complement,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.OptStatus,
		entity.Deleted,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_by_id.exec_error", "user_id", entity.ID, "error", execErr)
		return fmt.Errorf("update user by id: %w", execErr)
	}

	// Check if any rows were affected (does user exist?)
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.user.update_user_by_id.rows_affected_error", "user_id", entity.ID, "error", raErr)
		return fmt.Errorf("get rows affected: %w", raErr)
	}

	// Return sql.ErrNoRows if user not found (service layer maps to 404)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
