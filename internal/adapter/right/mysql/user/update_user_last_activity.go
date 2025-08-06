package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ua *UserAdapter) UpdateUserLastActivity(ctx context.Context, tx *sql.Tx, id int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `UPDATE users SET last_activity_at = ? WHERE id = ?;`

	_, err = ua.Update(ctx, tx, query,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserLastActivity: error executing Update", "error", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	return
}
