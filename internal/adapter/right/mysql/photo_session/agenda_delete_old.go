package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	photosessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/photo_session_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

var _ photosessionrepository.PhotoSessionRepositoryInterface = (*PhotoSessionAdapter)(nil)

// DeleteOldAgendaEntries removes agenda entries that ended before cutoff and are no longer tied to bookings.
func (a *PhotoSessionAdapter) DeleteOldAgendaEntries(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 {
		limit = 500
	}

	query := `DELETE FROM photographer_agenda_entries
        WHERE ends_at IS NOT NULL
          AND ends_at < ?
          AND NOT EXISTS (
              SELECT 1 FROM photographer_photo_session_bookings b
              WHERE b.agenda_entry_id = photographer_agenda_entries.id
          )
        LIMIT ?`

	res, execErr := a.ExecContext(ctx, tx, "delete_old_agenda_entries", query, cutoff, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.delete_old_agenda.exec_error", "cutoff", cutoff, "limit", limit, "error", execErr)
		return 0, fmt.Errorf("delete old agenda entries: %w", execErr)
	}

	rows, raErr := res.RowsAffected()
	if raErr != nil {
		logger.Warn("mysql.photo_session.delete_old_agenda.rows_affected_warning", "error", raErr)
		return 0, nil
	}

	if rows > 0 {
		logger.Debug("mysql.photo_session.delete_old_agenda.success", "deleted", rows, "cutoff", cutoff, "limit", limit)
	}

	return rows, nil
}
