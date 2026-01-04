package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListPhotographerIDsByLocation returns photographer IDs that cover the given city/state.
func (a *PhotoSessionAdapter) ListPhotographerIDsByLocation(ctx context.Context, tx *sql.Tx, city string, state string) ([]uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	trimmedCity := strings.TrimSpace(city)
	trimmedState := strings.TrimSpace(state)

	query := `
        SELECT DISTINCT u.id
        FROM users u
        JOIN user_roles ur ON ur.user_id = u.id
        JOIN roles r ON r.id = ur.role_id
        JOIN photographer_service_areas psa ON psa.photographer_user_id = u.id
        WHERE r.slug = 'photographer'
          AND r.is_active = 1
          AND ur.is_active = 1
          AND u.deleted = 0
          AND psa.city = ?
          AND psa.state = ?
    `

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, trimmedCity, trimmedState)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.photo_session.service_area.list_ids_location.query_error", "city", trimmedCity, "state", trimmedState, "err", queryErr)
		return nil, fmt.Errorf("list photographer ids by location: %w", queryErr)
	}
	defer rows.Close()

	ids := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		if scanErr := rows.Scan(&id); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.service_area.list_ids_location.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan photographer id by location: %w", scanErr)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.list_ids_location.rows_error", "err", err)
		return nil, fmt.Errorf("iterate photographer ids by location: %w", err)
	}

	return ids, nil
}
