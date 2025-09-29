package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// DeleteSessionsByUserID permanently removes all sessions for a given user
func (sa *SessionAdapter) DeleteSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	query := `DELETE FROM sessions WHERE user_id = ?`

	// Use helper when tx is provided
	if tx != nil {
		if _, err := sa.Delete(ctx, tx, query, userID); err != nil {
			slog.Error("sessionmysqladapter/DeleteSessionsByUserID: error executing Delete with tx", "error", err)
			return fmt.Errorf("delete sessions by user id: %w", err)
		}
		return nil
	}

	// Fallback if no transaction (should rarely happen in our flows)
	res, err := sa.db.DB.ExecContext(ctx, query, userID)
	if err != nil {
		slog.Error("sessionmysqladapter/DeleteSessionsByUserID: error executing ExecContext", "error", err)
		return fmt.Errorf("delete sessions by user id: %w", err)
	}
	if _, err := res.RowsAffected(); err != nil {
		slog.Warn("sessionmysqladapter/DeleteSessionsByUserID: failed to read rows affected", "error", err)
	}
	return nil
}
