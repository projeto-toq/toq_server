package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteVerticalComplexSize deletes a size row; returns sql.ErrNoRows when no row is affected.
func (a *PropertyCoverageAdapter) DeleteVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	const query = "DELETE FROM vertical_complex_sizes WHERE id = ? LIMIT 1;"
	return a.execUpdate(ctx, tx, "delete", query, id)
}
