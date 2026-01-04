package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateEntries inserts multiple agenda rows, preserving input order in the returned IDs. No-op when entries is empty.
// Expects a non-nil transaction when atomicity across multiple inserts is required; traces and logs infra errors.
func (a *PhotoSessionAdapter) CreateEntries(ctx context.Context, tx *sql.Tx, entries []photosessionmodel.AgendaEntryInterface) ([]uint64, error) {
	if len(entries) == 0 {
		return nil, nil
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO photographer_agenda_entries (
		photographer_user_id, entry_type, source, source_id, starts_at, ends_at, blocking, reason, timezone
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	ids := make([]uint64, 0, len(entries))
	for _, entry := range entries {
		entity := converters.ToAgendaEntryEntity(entry)

		var sourceID any
		if entity.SourceID.Valid {
			sourceID = entity.SourceID.Int64
		}

		var reason any
		if entity.Reason.Valid {
			reason = entity.Reason.String
		}

		result, execErr := a.ExecContext(
			ctx,
			tx,
			"insert",
			query,
			entity.PhotographerUserID,
			entity.EntryType,
			entity.Source,
			sourceID,
			entity.StartsAt,
			entity.EndsAt,
			entity.Blocking,
			reason,
			entity.Timezone,
		)
		if execErr != nil {
			utils.SetSpanError(ctx, execErr)
			logger.Error("mysql.photo_session.create_entries.exec_error", "photographer_id", entity.PhotographerUserID, "err", execErr)
			return nil, fmt.Errorf("insert photographer agenda entry: %w", execErr)
		}

		id, lastErr := result.LastInsertId()
		if lastErr != nil {
			utils.SetSpanError(ctx, lastErr)
			logger.Error("mysql.photo_session.create_entries.last_id_error", "photographer_id", entity.PhotographerUserID, "err", lastErr)
			return nil, fmt.Errorf("agenda entry last insert id: %w", lastErr)
		}

		entry.SetID(uint64(id))
		ids = append(ids, uint64(id))
	}

	return ids, nil
}
