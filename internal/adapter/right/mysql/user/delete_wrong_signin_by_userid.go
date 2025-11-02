package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteWrongSignInByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM temp_wrong_signin WHERE user_id = ?;`

	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_wrong_signin.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete temp_wrong_signin: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_wrong_signin.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete temp_wrong_signin rows affected: %w", rowsErr)
	}

	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return rowsAffected, nil
}
