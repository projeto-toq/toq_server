package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM user_roles WHERE user_id = ?;`

	deleted, err = ua.Delete(ctx, tx, query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.delete_user_roles.delete_error", "error", err)
		return 0, fmt.Errorf("delete user_roles by user_id: %w", err)
	}

	if deleted == 0 {
		errNoRows := errors.New("no user_roles rows deleted")
		utils.SetSpanError(ctx, errNoRows)
		logger.Error("mysql.user.delete_user_roles.no_rows_deleted", "user_id", userID, "error", errNoRows)
		return 0, errNoRows
	}

	return
}
