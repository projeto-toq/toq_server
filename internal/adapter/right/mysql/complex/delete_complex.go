package mysqlcomplexadapter

import (
	"context"
	"database/sql"
)

func (ca *ComplexAdapter) DeleteComplex(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	query := "DELETE FROM complex WHERE id = ? LIMIT 1;"
	return ca.Delete(ctx, tx, query, id)
}
