package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Read executa uma query SELECT e retorna os resultados
func (pa *PermissionAdapter) Read(ctx context.Context, tx *sql.Tx, query string, args ...any) ([][]any, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqlpermissionadapter/Read: error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		slog.Error("mysqlpermissionadapter/Read: error executing query", "error", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		slog.Error("mysqlpermissionadapter/Read: error getting columns", "error", err)
		return nil, err
	}

	var results [][]any
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			slog.Error("mysqlpermissionadapter/Read: error scanning row", "error", err)
			return nil, err
		}

		results = append(results, values)
	}

	if err := rows.Err(); err != nil {
		slog.Error("mysqlpermissionadapter/Read: error iterating rows", "error", err)
		return nil, err
	}

	return results, nil
}
