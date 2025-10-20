package mysqlcomplexadapter

import (
	"context"
	"database/sql"
)

func (ca *ComplexAdapter) DeleteComplexSize(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	query := "DELETE FROM complex_sizes WHERE id = ? LIMIT 1;"
	return ca.Delete(ctx, tx, query, id)
}
