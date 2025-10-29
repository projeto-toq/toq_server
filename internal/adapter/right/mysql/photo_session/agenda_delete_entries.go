package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) DeleteEntriesBySource(ctx context.Context, tx *sql.Tx, photographerID uint64, entryType photosessionmodel.AgendaEntryType, source photosessionmodel.AgendaEntrySource, sourceID *uint64) (int64, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return 0, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM photographer_agenda_entries WHERE photographer_user_id = ? AND entry_type = ? AND source = ?`
	args := []any{photographerID, string(entryType), string(source)}

	if sourceID != nil {
		query += ` AND source_id = ?`
		args = append(args, *sourceID)
	} else {
		query += ` AND source_id IS NULL`
	}

	result, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.delete_entries.exec_error", "photographer_id", photographerID, "err", err)
		return 0, fmt.Errorf("delete agenda entries by source: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.delete_entries.rows_error", "photographer_id", photographerID, "err", err)
		return 0, fmt.Errorf("rows affected agenda entries: %w", err)
	}

	return affected, nil
}
