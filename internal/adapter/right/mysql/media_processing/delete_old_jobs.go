package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	mediaprocrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/media_processing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

var _ mediaprocrepository.RepositoryInterface = (*MediaProcessingAdapter)(nil)

// DeleteOldJobs removes terminal media processing jobs older than the cutoff (by completed_at or created_at fallback).
// Returns rows deleted; zero rows means nothing to prune.
func (a *MediaProcessingAdapter) DeleteOldJobs(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error) {
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

	query := `DELETE FROM media_processing_jobs
        WHERE status IN ('SUCCEEDED','PARTIAL_SUCCESS','FAILED')
          AND COALESCE(completed_at, created_at) < ?
        LIMIT ?`

	res, execErr := a.ExecContext(ctx, tx, "delete_old_media_jobs", query, cutoff, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.media_processing.delete_old_jobs.exec_error", "cutoff", cutoff, "limit", limit, "error", execErr)
		return 0, fmt.Errorf("delete old media processing jobs: %w", execErr)
	}

	rows, raErr := res.RowsAffected()
	if raErr != nil {
		logger.Warn("mysql.media_processing.delete_old_jobs.rows_affected_warning", "error", raErr)
		return 0, nil
	}

	if rows > 0 {
		logger.Debug("mysql.media_processing.delete_old_jobs.success", "deleted", rows, "cutoff", cutoff, "limit", limit)
	}

	return rows, nil
}
