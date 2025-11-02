package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) DeleteEntriesBySource(ctx context.Context, tx *sql.Tx, photographerID uint64, entryType photosessionmodel.AgendaEntryType, source photosessionmodel.AgendaEntrySource, sourceID *uint64) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()
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

	result, execErr := a.ExecContext(ctx, tx, "delete", query, args...)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.delete_entries.exec_error", "photographer_id", photographerID, "err", execErr)
		return 0, fmt.Errorf("delete agenda entries by source: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.delete_entries.rows_error", "photographer_id", photographerID, "err", rowsErr)
		return 0, fmt.Errorf("rows affected agenda entries: %w", rowsErr)
	}

	return affected, nil
}
