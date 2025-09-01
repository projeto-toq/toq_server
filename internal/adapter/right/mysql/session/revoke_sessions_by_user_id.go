package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) RevokeSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	query := `UPDATE sessions SET revoked = true WHERE user_id = ? AND revoked = false`

	_, err = sa.Update(ctx, tx, query, userID)
	if err != nil {
		slog.Error("sessionmysqladapter/RevokeSessionsByUserID: error executing Update", "error", err)
		return fmt.Errorf("revoke sessions by user id: %w", err)
	}

	return nil
}
