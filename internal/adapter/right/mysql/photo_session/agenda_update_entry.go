package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateEntry persists changes applied to an agenda entry.
func (a *PhotoSessionAdapter) UpdateEntry(ctx context.Context, tx *sql.Tx, entry photosessionmodel.AgendaEntryInterface) error {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToAgendaEntryEntity(entry)

	query := `UPDATE photographer_agenda_entries
        SET entry_type = ?, source = ?, source_id = ?, starts_at = ?, ends_at = ?, blocking = ?, reason = ?, timezone = ?, updated_at = NOW()
        WHERE id = ?`

	var sourceID any
	if entity.SourceID.Valid {
		sourceID = entity.SourceID.Int64
	}

	var reason any
	if entity.Reason.Valid {
		reason = entity.Reason.String
	}

	if _, err := exec.ExecContext(
		ctx,
		query,
		entity.EntryType,
		entity.Source,
		sourceID,
		entity.StartsAt,
		entity.EndsAt,
		entity.Blocking,
		reason,
		entity.Timezone,
		entity.ID,
	); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_entry.exec_error", "entry_id", entity.ID, "err", err)
		return fmt.Errorf("update agenda entry: %w", err)
	}

	return nil
}
