package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Read executa uma query SELECT e retorna os resultados
func (pa *PermissionAdapter) Read(ctx context.Context, tx *sql.Tx, query string, args ...any) ([][]any, error) {
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
		logger.Error("mysql.permission.read.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare permission read statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.read.query_error", "error", err)
		return nil, fmt.Errorf("query permission read: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.read.columns_error", "error", err)
		return nil, fmt.Errorf("columns permission read: %w", err)
	}

	var results [][]any
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.permission.read.scan_error", "error", err)
			return nil, fmt.Errorf("scan permission read row: %w", err)
		}

		results = append(results, values)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.read.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration permission read: %w", err)
	}

	return results, nil
}

// ReadRow removido: para leituras de linha única, use QueryRowContext + Scan tipado no método específico.
