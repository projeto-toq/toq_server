package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListAllUsers retrieves all non-deleted users from the database
//
// This method implements the UserRepoPortInterface.ListAllUsers() contract.
// Returns all users where deleted=0 without pagination or filtering.
//
// Naming Convention: Follows Section 8.1.4 of guide - List* prefix for collection retrieval
//
// Parameters:
//   - ctx: Context for tracing, logging, cancellation
//   - tx: Optional transaction (*sql.Tx can be nil for auto-commit)
//
// Returns:
//   - users: Slice of UserInterface domain models
//   - err: sql.ErrNoRows if no users found, or database/conversion errors
//
// Query Details:
//   - Table: users
//   - Filter: deleted = 0 (only active users)
//   - Columns: All 22 user columns explicitly listed (no SELECT *)
//   - Order: Natural order (insertion order, no explicit ORDER BY)
//
// Observability:
//   - OpenTelemetry span: "ListAllUsers"
//   - Structured logs: mysql.user.list_all_users.*
//   - Metrics: Query duration, row count via InstrumentedAdapter
func (ua *UserAdapter) ListAllUsers(ctx context.Context, tx *sql.Tx) (users []usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query all active users (deleted=0)
	// Explicitly lists all 22 columns (no SELECT * per Section 8.3.2)
	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE deleted = 0`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.list_all_users.query_error", "error", queryErr)
		return nil, fmt.Errorf("list all users query: %w", queryErr)
	}
	defer rows.Close()

	entities, err := scanUserEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.list_all_users.scan_error", "error", err)
		return nil, fmt.Errorf("scan user rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, entity := range entities {
		user := userconverters.UserEntityToDomain(entity)
		users = append(users, user)
	}

	return users, nil
}
