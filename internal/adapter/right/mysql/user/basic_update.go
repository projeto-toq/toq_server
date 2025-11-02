package mysqluseradapter

import (
	"context"
	"database/sql"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) Update(ctx context.Context, tx *sql.Tx, query string, args ...any) (affected int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	startedAt := time.Now()
	defer func() {
		ua.Observe("update", query, time.Since(startedAt))
	}()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.update.prepare_error", "error", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.update.exec_error", "error", err)
		return 0, err
	}

	affected, err = result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.update.rows_affected_error", "error", err)
		return 0, err
	}

	return affected, nil
}
