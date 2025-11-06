package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserByID retrieves a user by their unique ID from the users table.
// Returns sql.ErrNoRows if no user is found with the given ID or if the user is marked as deleted.
// This function ensures only active (non-deleted) users are returned.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - id: User's unique identifier
//
// Returns:
//   - user: UserInterface containing all user data
//   - error: sql.ErrNoRows if not found, or other database errors
//
// Business Rules:
//   - Query filters by deleted = 0 (soft delete pattern)
//   - Returns error if multiple users found with same ID (data integrity check)
//
// Note: Active role is NOT populated by this function. Caller must use
//
//	Permission Service to set the active role after retrieval.
func (ua *UserAdapter) GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error) {
	// Initialize tracing for observability (metrics + distributed tracing)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query only active users (deleted = 0) to maintain data integrity
	// Note: photo column removed - does not exist in current schema
	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, 
	          creci_validity, born_at, phone_number, email, zip_code, street, number, 
	          complement, neighborhood, city, state, password, opt_status, 
	          last_activity_at, deleted, last_signin_attempt 
	          FROM users 
	          WHERE id = ? AND deleted = 0`

	// Execute query using instrumented adapter (auto-generates metrics + tracing)
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, id)
	if queryErr != nil {
		// Mark span as error for distributed tracing analysis
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_by_id.query_error", "user_id", id, "error", queryErr)
		return nil, fmt.Errorf("query user by id: %w", queryErr)
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entities (type-safe scanning)
	entities, err := scanUserEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_id.scan_error", "user_id", id, "error", err)
		return nil, fmt.Errorf("scan user rows: %w", err)
	}

	// Handle no results: return standard sql.ErrNoRows for service layer handling
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: unique constraint should prevent multiple rows, but verify
	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple users found with the same ID: %d", id)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_id.multiple_users_error", "user_id", id, "count", len(entities), "error", errMultiple)
		return nil, errMultiple
	}

	// Convert database entity to domain model (separation of concerns)
	user = userconverters.UserEntityToDomain(entities[0])

	return user, nil
}
