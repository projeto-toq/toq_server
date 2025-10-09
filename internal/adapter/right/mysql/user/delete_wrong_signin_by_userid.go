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

	deleted, err = ua.Delete(ctx, tx, query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.delete_wrong_signin.delete_error", "error", err)
		return 0, fmt.Errorf("delete temp_wrong_signin: %w", err)
	}

	if deleted == 0 {
		return 0, sql.ErrNoRows
	}

	return
}
