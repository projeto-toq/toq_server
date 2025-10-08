package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	deleted, err = ua.Delete(ctx, tx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.delete_invite.delete_error", "error", err)
		return 0, fmt.Errorf("delete invite by id: %w", err)
	}

	if deleted == 0 {
		errNoRows := errors.New("no agency_invites rows deleted")
		utils.SetSpanError(ctx, errNoRows)
		logger.Error("mysql.user.delete_invite.no_rows_deleted", "invite_id", id, "error", errNoRows)
		return 0, errNoRows
	}

	return
}
