package mysqluseradapter

import (
	"context"
	"database/sql"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) Create(ctx context.Context, tx *sql.Tx, query string, args ...any) (id int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	startedAt := time.Now()
	defer func() {
		ua.Observe("insert", query, time.Since(startedAt))
	}()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.create.prepare_error", "error", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.create.exec_error", "error", err)
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.create.last_insert_error", "error", err)
		return 0, err
	}

	return id, nil
}
