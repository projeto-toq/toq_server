package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateUserPasswordByID updates only the password hash for a specific user
//
// This function performs a targeted update of the password field without modifying
// other user data. Used for password reset and password change operations.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - user: UserInterface with ID and new password hash set
//
// Returns:
//   - error: sql.ErrNoRows if user not found, database errors
//
// Business Rules:
//   - Password must be bcrypt hashed BEFORE calling this function
//   - Updates only password field (all other fields remain unchanged)
//   - Does NOT validate old password (service layer responsibility)
//   - Does NOT check deleted status (updates any user, use with caution)
//
// Security Considerations:
//   - NEVER store plain text passwords
//   - Password must be hashed with bcrypt by service layer
//   - Old password validation performed by service before calling this
//   - Consider logging password change event in audit table
//
// Edge Cases:
//   - User deleted: Password updated but user invisible (consider checking deleted=0)
//   - User ID 0: Invalid, but UPDATE will succeed with 0 rows affected
//   - Empty password hash: Allowed by DB but should be prevented by service
//
// Performance:
//   - Single-row update using PRIMARY KEY (very fast)
//   - No indexes impacted (password field not indexed)
//
// Example:
//
//	user := usermodel.NewUser()
//	user.SetID(userID)
//	user.SetPassword(bcryptHash) // Pre-hashed by service
//
//	err := adapter.UpdateUserPasswordByID(ctx, tx, user)
//	if err != nil {
//	    // Handle error (rollback transaction in service)
//	}
func (ua *UserAdapter) UpdateUserPasswordByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update only password field by user ID
	// Note: No WHERE deleted = 0 check - updates ANY user (active or deleted)
	query := `UPDATE users SET password = ? WHERE id = ?;`

	// Execute update using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		user.GetPassword(),
		user.GetID(),
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_password.exec_error", "user_id", user.GetID(), "error", execErr)
		return fmt.Errorf("update user password: %w", execErr)
	}

	// Check rows affected (validate user exists)
	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_password.rows_affected_error", "user_id", user.GetID(), "error", rowsErr)
		return fmt.Errorf("user password update rows affected: %w", rowsErr)
	}

	return
}
