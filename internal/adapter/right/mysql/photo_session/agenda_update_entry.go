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
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToAgendaEntryEntity(entry)

	query := `UPDATE photographer_agenda_entries
        SET entry_type = ?, source = ?, source_id = ?, starts_at = ?, ends_at = ?, blocking = ?, reason = ?, timezone = ?
        WHERE id = ?`

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
		"update",
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
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.update_entry.exec_error", "entry_id", entity.ID, "err", execErr)
		return fmt.Errorf("update agenda entry: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.update_entry.rows_error", "entry_id", entity.ID, "err", rowsErr)
		return fmt.Errorf("rows affected agenda entry: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
