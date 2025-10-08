package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) Read(ctx context.Context, tx *sql.Tx, query string, args ...any) (entity [][]any, err error) {
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
		logger.Error("mysql.complex.read.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare complex read statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.read.query_error", "error", err)
		return nil, fmt.Errorf("query complex read: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.read.columns_error", "error", err)
		return nil, fmt.Errorf("columns complex read: %w", err)
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
			logger.Error("mysql.complex.read.scan_error", "error", err)
			return nil, fmt.Errorf("scan complex read row: %w", err)
		}
		entity = append(entity, entityElements)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.read.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration complex read: %w", err)
	}

	return entity, nil
}
