package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListPhotographerIDs returns distinct photographer IDs that have agenda entries.
func (a *PhotoSessionAdapter) ListPhotographerIDs(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT DISTINCT photographer_user_id FROM photographer_agenda_entries`

	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_photographers.query_error", "err", err)
		return nil, fmt.Errorf("list photographer ids: %w", err)
	}
	defer rows.Close()

	ids := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		if scanErr := rows.Scan(&id); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.list_photographers.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan photographer id: %w", scanErr)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_photographers.rows_error", "err", err)
		return nil, fmt.Errorf("iterate photographer ids: %w", err)
	}

	return ids, nil
}
