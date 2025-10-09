package mysqlglobaladapter

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ga *GlobalAdapter) Update(ctx context.Context, tx *sql.Tx, query string, args ...any) (affected int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.common.update.prepare_error", "error", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.common.update.exec_error", "error", err)
		return 0, err
	}

	affected, err = result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.common.update.rows_affected_error", "error", err)
		return 0, err
	}

	return affected, nil
}
