package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserPasswordByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE users SET password = ? WHERE id = ?;`

	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		user.GetPassword(),
		user.GetID(),
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_password.exec_error", "error", execErr)
		return fmt.Errorf("update user password: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_password.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user password update rows affected: %w", rowsErr)
	}

	return
}
