package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserByNationalID retrieves a user by their national ID (CPF or CNPJ)
//
// This function searches for active (non-deleted) users by their unique national identifier.
// Returns sql.ErrNoRows if no matching user is found.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - nationalID: User's CPF (11 digits) or CNPJ (14 digits) without formatting
//
// Returns:
//   - user: UserInterface with all user data
//   - error: sql.ErrNoRows if not found, or database errors
//
// Business Rules:
//   - National ID is UNIQUE constraint in database
//   - Query does NOT filter by deleted (returns even deleted users for security checks)
//
// Security Note:
//   - Used for authentication and duplicate detection
//   - Sensitive PII - ensure proper logging controls
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

	// Query by national_id (unique constraint) - includes deleted users for security checks
	// Note: photo column removed - does not exist in current schema
	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, 
	          born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, 
	          password, opt_status, last_activity_at, deleted, last_signin_attempt 
	          FROM users WHERE national_id = ?`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, nationalID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_by_national_id.query_error", "error", queryErr)
		return nil, fmt.Errorf("get user by national_id query: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entities
	entities, err := scanUserEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_national_id.scan_error", "error", err)
		return nil, fmt.Errorf("scan user by national id rows: %w", err)
	}

	// Handle no results
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: unique constraint should prevent multiple rows
	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple users found for national_id: %s", nationalID)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_national_id.multiple_users_error", "national_id", nationalID, "count", len(entities), "error", errMultiple)
		return nil, errMultiple
	}

	// Convert database entity to domain model
	user = userconverters.UserEntityToDomain(entities[0])

	return user, nil
}
