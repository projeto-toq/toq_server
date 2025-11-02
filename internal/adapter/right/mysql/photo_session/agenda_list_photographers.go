package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListPhotographerIDs returns distinct photographer IDs that are registered with the photographer role.
func (a *PhotoSessionAdapter) ListPhotographerIDs(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT DISTINCT u.id
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE r.slug = 'photographer'
		  AND r.is_active = 1
		  AND ur.is_active = 1
		  AND u.deleted = 0
	`

	rows, queryErr := a.QueryContext(ctx, tx, "select", query)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.photo_session.list_photographers.query_error", "err", queryErr)
		return nil, fmt.Errorf("list photographer ids: %w", queryErr)
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
