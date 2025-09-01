package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `DELETE FROM user_roles WHERE user_id = ?;`

	deleted, err = ua.Delete(ctx, tx, query, userID)
	if err != nil {
		slog.Error("mysqluseradapter/DeleteUserRolesByUserID: error executing Delete", "error", err)
		return 0, fmt.Errorf("delete user_roles by user_id: %w", err)
	}

	if deleted == 0 {
		return 0, errors.New("no user_roles rows deleted")
	}

	return
}
