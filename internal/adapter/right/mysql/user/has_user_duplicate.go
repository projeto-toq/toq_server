package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HasUserDuplicate checks if any active user exists with the given phone, email, or national ID
//
// This function is used during registration to prevent duplicate accounts.
// Searches for any non-deleted user matching at least one of the unique identifiers.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction
//   - user: UserInterface with phone, email, and national ID to check
//
// Returns:
//   - exists: true if any active user found with matching identifiers
//   - error: Database errors
//
// Business Rules:
//   - Only checks non-deleted users (deleted = 0)
//   - Uses OR logic: match on ANY of the three identifiers triggers duplicate
func (ua *UserAdapter) HasUserDuplicate(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query checks for any active user with matching phone, email, or national ID
	query := `SELECT count(id) as count
				FROM users WHERE (phone_number = ? OR email = ? OR national_id = ? ) AND deleted = 0`

	row := ua.QueryRowContext(ctx, tx, "select", query,
		user.GetPhoneNumber(),
		user.GetEmail(),
		user.GetNationalID(),
	)

	var qty int64
	if scanErr := row.Scan(&qty); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.has_duplicate.scan_error", "error", scanErr)
		return false, fmt.Errorf("check user duplicate: %w", scanErr)
	}

	exist = qty > 0

	return exist, nil
}
