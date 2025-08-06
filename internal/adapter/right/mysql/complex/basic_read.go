package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) Read(ctx context.Context, tx *sql.Tx, query string, args ...any) (entity [][]any, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqlcomplexadapter: error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		slog.Error("mysqlcomplexadapter: error executing query", "error", err)
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		slog.Error("mysqlcomplexadapter: error fetching columns", "error", err)
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
			slog.Error("mysqlcomplexadapter: error scanning row", "error", err)
			return nil, err
		}
		entity = append(entity, entityElements)
	}

	if err = rows.Err(); err != nil {
		slog.Error("mysqlcomplexadapter: rows iteration error", "error", err)
		return nil, err
	}

	return entity, nil
}
