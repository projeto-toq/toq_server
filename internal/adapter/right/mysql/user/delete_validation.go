package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteValidation(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM temp_user_validations WHERE user_id = ?;`

	result, execErr := ua.ExecContext(ctx, tx, "delete", query, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_validation.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete validation: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_validation.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete validation rows affected: %w", rowsErr)
	}

	// Idempotent: if no rows were deleted, that's fine (nothing to clean up)
	return rowsAffected, nil
}
