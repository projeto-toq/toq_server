package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sa *SessionAdapter) MarkSessionRotated(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	query := `UPDATE sessions SET rotated_at = UTC_TIMESTAMP() WHERE id = ?`

	_, err = sa.Update(ctx, tx, query, id)
	if err != nil {
		slog.Error("sessionmysqladapter/MarkSessionRotated: error executing Update", "error", err)
		return status.Error(codes.Internal, "Failed to mark session as rotated")
	}

	return nil
}
