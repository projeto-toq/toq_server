package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteInviteByID(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM agency_invites WHERE id = ?;`

	result, execErr := ua.ExecContext(ctx, tx, "delete", query, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_invite.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete invite by id: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_invite.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete invite rows affected: %w", rowsErr)
	}

	if rowsAffected == 0 {
		errNoRows := errors.New("no agency_invites rows deleted")
		utils.SetSpanError(ctx, errNoRows)
		logger.Error("mysql.user.delete_invite.no_rows_deleted", "invite_id", id, "error", errNoRows)
		return 0, errNoRows
	}

	return rowsAffected, nil
}
