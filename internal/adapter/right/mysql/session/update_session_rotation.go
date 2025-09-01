package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) UpdateSessionRotation(ctx context.Context, tx *sql.Tx, id int64, rotationCounter int, lastRefreshAt time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	query := `UPDATE sessions SET rotation_counter = ?, last_refresh_at = ? WHERE id = ?`

	_, err = sa.Update(ctx, tx, query, rotationCounter, lastRefreshAt, id)
	if err != nil {
		slog.Error("sessionmysqladapter/UpdateSessionRotation: error executing Update", "error", err)
		return fmt.Errorf("update session rotation: %w", err)
	}

	return nil
}
