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

func (a *PhotoSessionAdapter) ListEntriesByRange(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.AgendaEntryInterface, error) {
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

	query := `SELECT id, photographer_user_id, entry_type, source, source_id, starts_at, ends_at, blocking, reason, timezone, created_at, updated_at
		FROM photographer_agenda_entries
		WHERE photographer_user_id = ? AND ends_at > ? AND starts_at < ?
		ORDER BY starts_at ASC`

	rows, err := exec.QueryContext(ctx, query, photographerID, rangeStart, rangeEnd)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_entries.query_error", "photographer_id", photographerID, "err", err)
		return nil, fmt.Errorf("list photographer agenda entries: %w", err)
	}
	defer rows.Close()

	result := make([]photosessionmodel.AgendaEntryInterface, 0)
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
			&row.CreatedAt,
			&row.UpdatedAt,
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.list_entries.scan_error", "photographer_id", photographerID, "err", scanErr)
			return nil, fmt.Errorf("scan photographer agenda entry: %w", scanErr)
		}

		model := converters.ToAgendaEntryModel(row)
		result = append(result, model)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_entries.rows_error", "photographer_id", photographerID, "err", err)
		return nil, fmt.Errorf("iterate photographer agenda entries: %w", err)
	}

	return result, nil
}
