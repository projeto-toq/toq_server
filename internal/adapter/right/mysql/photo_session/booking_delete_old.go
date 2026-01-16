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

// DeleteOldBookings removes bookings in terminal states whose ends_at is older than cutoff.
func (a *PhotoSessionAdapter) DeleteOldBookings(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error) {
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

	query := `DELETE FROM photographer_photo_session_bookings
        WHERE ends_at IS NOT NULL
          AND ends_at < ?
          AND status IN ('CANCELLED','REJECTED','DONE')
        LIMIT ?`

	res, execErr := a.ExecContext(ctx, tx, "delete_old_bookings", query, cutoff, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.delete_old_bookings.exec_error", "cutoff", cutoff, "limit", limit, "error", execErr)
		return 0, fmt.Errorf("delete old bookings: %w", execErr)
	}

	rows, raErr := res.RowsAffected()
	if raErr != nil {
		logger.Warn("mysql.photo_session.delete_old_bookings.rows_affected_warning", "error", raErr)
		return 0, nil
	}

	if rows > 0 {
		logger.Debug("mysql.photo_session.delete_old_bookings.success", "deleted", rows, "cutoff", cutoff, "limit", limit)
	}
	return rows, nil
}
