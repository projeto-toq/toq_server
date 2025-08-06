package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ua *UserAdapter) DeleteInviteByID(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `DELETE FROM agency_invites WHERE id = ?;`

	deleted, err = ua.Delete(ctx, tx, query, id)
	if err != nil {
		slog.Error("mysqluseradapter/DeleteInviteByID: error executing Delete", "error", err)
		return 0, status.Error(codes.Internal, "internal server error")
	}

	if deleted == 0 {
		return 0, status.Error(codes.NotFound, "agency invite not found")
	}

	return
}
