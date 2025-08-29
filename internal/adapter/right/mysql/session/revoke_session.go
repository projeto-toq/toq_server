package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
		return utils.ErrInternalServer
	}

	return nil
}
