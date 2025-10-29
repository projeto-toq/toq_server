package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) CreateEntries(ctx context.Context, tx *sql.Tx, entries []photosessionmodel.AgendaEntryInterface) ([]uint64, error) {
	if len(entries) == 0 {
		return nil, nil
	}

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

		result, err := exec.ExecContext(
			ctx,
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
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.create_entries.exec_error", "photographer_id", entity.PhotographerUserID, "err", err)
			return nil, fmt.Errorf("insert photographer agenda entry: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.create_entries.last_id_error", "photographer_id", entity.PhotographerUserID, "err", err)
			return nil, fmt.Errorf("agenda entry last insert id: %w", err)
		}

		entry.SetID(uint64(id))
		ids = append(ids, uint64(id))
	}

	return ids, nil
}
