package mysqlglobaladapter

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ga *GlobalAdapter) Read(ctx context.Context, tx *sql.Tx, query string, args ...any) (entity [][]any, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.common.read.prepare_error", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.common.read.query_error", "error", err)
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.common.read.columns_error", "error", err)
		return nil, err
	}

	entity = make([][]any, 0)

	for rows.Next() {
		entityElements := make([]any, len(cols))
		entityElementPtrs := make([]any, len(cols))
		for i := range entityElements {
			entityElementPtrs[i] = &entityElements[i]
		}
		err = rows.Scan(entityElementPtrs...)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.common.read.scan_error", "error", err)
			return nil, err
		}
		entity = append(entity, entityElements)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.common.read.rows_error", "error", err)
		return nil, err
	}

	return entity, nil
}
