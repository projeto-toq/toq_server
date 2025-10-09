package mysqluseradapter

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) Delete(ctx context.Context, tx *sql.Tx, query string, args ...any) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.delete.prepare_error", "error", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.delete.exec_error", "error", err)
		return
	}

	deleted, err = result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.delete.rows_affected_error", "error", err)
		return
	}

	return
}
