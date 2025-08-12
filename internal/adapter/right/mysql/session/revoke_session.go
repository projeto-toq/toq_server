package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sa *SessionAdapter) RevokeSession(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	query := `UPDATE sessions SET revoked = true WHERE id = ?`

	_, err = sa.Update(ctx, tx, query, id)
	if err != nil {
		slog.Error("sessionmysqladapter/RevokeSession: error executing Update", "error", err)
		return status.Error(codes.Internal, "Failed to revoke session")
	}

	return nil
}
