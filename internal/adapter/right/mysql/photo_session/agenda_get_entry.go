package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) GetEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) (photosessionmodel.AgendaEntryInterface, error) {
	return a.getEntry(ctx, tx, entryID, false)
}

func (a *PhotoSessionAdapter) GetEntryByIDForUpdate(ctx context.Context, tx *sql.Tx, entryID uint64) (photosessionmodel.AgendaEntryInterface, error) {
	return a.getEntry(ctx, tx, entryID, true)
}

func (a *PhotoSessionAdapter) getEntry(ctx context.Context, tx *sql.Tx, entryID uint64, forUpdate bool) (photosessionmodel.AgendaEntryInterface, error) {
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

	query := `SELECT id, photographer_user_id, entry_type, source, source_id, starts_at, ends_at, blocking, reason, timezone
		FROM photographer_agenda_entries WHERE id = ?`
	if forUpdate {
		query += " FOR UPDATE"
	}

	row := entity.AgendaEntry{}
	scanErr := exec.QueryRowContext(ctx, query, entryID).Scan(
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
)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, scanErr
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.photo_session.get_entry.scan_error", "entry_id", entryID, "err", scanErr)
		return nil, fmt.Errorf("get agenda entry: %w", scanErr)
	}

	return converters.ToAgendaEntryModel(row), nil
}
