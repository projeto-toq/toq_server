package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateUser inserts a new user record into the users table
//
// This function creates a new user with all provided data and populates the
// user's ID field with the auto-generated primary key.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - user creation must be transactional)
//   - user: UserInterface with all required fields populated (ID will be ignored)
//
// Returns:
//   - error: Database errors (constraint violations, connection issues)
//
// Side Effects:
//   - Modifies user object by setting ID to the newly inserted row's primary key
//
// Business Rules:
//   - All mandatory fields must be set (enforced by NOT NULL constraints)
//   - Email, phone, and national ID must be unique (UNIQUE constraints)
//   - Password must be bcrypt hash (enforced by service layer, not validated here)
//
// Unique Constraint Violations:
//   - Duplicate email: error contains "Duplicate entry" and "email"
//   - Duplicate phone: error contains "Duplicate entry" and "phone_number"
//   - Duplicate CPF/CNPJ: error contains "Duplicate entry" and "national_id"
//
// Note: Service layer maps constraint violations to domain errors (409 Conflict)
func (ua *UserAdapter) CreateUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// SQL insert with all user fields
	sql := `INSERT INTO users (
			full_name, nick_name, national_id, creci_number, creci_state, creci_validity,
			born_at, phone_number, email, zip_code, street, number, complement, neighborhood, 
			city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Convert domain model to database entity
	entity := userconverters.UserDomainToEntity(user)

	// Execute insert using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "insert", sql,
		entity.FullName,
		entity.NickName,
		entity.NationalID,
		entity.CreciNumber,
		entity.CreciState,
		entity.CreciValidity,
		entity.BornAT,
		entity.PhoneNumber,
		entity.Email,
		entity.ZipCode,
		entity.Street,
		entity.Number,
		entity.Complement,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.Password,
		entity.OptStatus,
		entity.LastActivityAT,
		entity.Deleted,
		entity.LastSignInAttempt,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.create_user.exec_error", "error", execErr)
		return fmt.Errorf("create user: %w", execErr)
	}

	// Retrieve auto-generated primary key
	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.user.create_user.last_insert_id_error", "error", lastErr)
		return fmt.Errorf("user last insert id: %w", lastErr)
	}

	// Update user object with generated ID (side effect)
	user.SetID(id)

	return
}
