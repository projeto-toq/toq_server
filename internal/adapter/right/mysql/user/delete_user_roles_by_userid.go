package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return 0, status.Error(codes.Internal, "internal server error")
	}

	if deleted == 0 {
		return 0, status.Error(codes.NotFound, "user roles not found")
	}

	return
}
