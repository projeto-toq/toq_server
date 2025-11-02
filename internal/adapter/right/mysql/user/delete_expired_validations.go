package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteExpiredValidations deletes rows from temp_user_validations where all codes are either empty/NULL or expired
// The deletion is limited by the given limit to avoid long-running operations
func (ua *UserAdapter) DeleteExpiredValidations(ctx context.Context, tx *sql.Tx, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// A row can be deleted when none of the three codes is currently valid (all are empty/NULL or expired)
	query := `DELETE FROM temp_user_validations
		WHERE
			( (email_code IS NULL OR email_code = '' OR email_code_exp < NOW())
			AND (phone_code IS NULL OR phone_code = '' OR phone_code_exp < NOW())
			AND (password_code IS NULL OR password_code = '' OR password_code_exp < NOW()) )
		LIMIT ?;`

	result, execErr := ua.ExecContext(ctx, tx, "delete", query, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_expired_validations.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete expired validations: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_expired_validations.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete expired validations rows affected: %w", rowsErr)
	}

	return rowsAffected, nil
}
