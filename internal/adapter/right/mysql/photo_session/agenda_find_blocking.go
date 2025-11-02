package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) FindBlockingEntries(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.AgendaEntryInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, photographer_user_id, entry_type, source, source_id, starts_at, ends_at, blocking, reason, timezone
		FROM photographer_agenda_entries
		WHERE photographer_user_id = ? AND blocking = 1 AND ends_at > ? AND starts_at < ?
		ORDER BY starts_at ASC`

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, photographerID, rangeStart, rangeEnd)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.photo_session.find_blocking.query_error", "photographer_id", photographerID, "err", queryErr)
		return nil, fmt.Errorf("find blocking agenda entries: %w", queryErr)
	}
	defer rows.Close()

	entries := make([]photosessionmodel.AgendaEntryInterface, 0)
	for rows.Next() {
		row := entity.AgendaEntry{}
		if scanErr := rows.Scan(
			&row.ID,
			&row.PhotographerUserID,
			&row.EntryType,
			&row.Source,
			&row.SourceID,
			&row.StartsAt,
			&row.EndsAt,
			&row.Blocking,
			&row.Reason,
			&row.Timezone,
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.find_blocking.scan_error", "photographer_id", photographerID, "err", scanErr)
			return nil, fmt.Errorf("scan blocking agenda entry: %w", scanErr)
		}
		entries = append(entries, converters.ToAgendaEntryModel(row))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.find_blocking.rows_error", "photographer_id", photographerID, "err", err)
		return nil, fmt.Errorf("iterate blocking agenda entries: %w", err)
	}

	return entries, nil
}
