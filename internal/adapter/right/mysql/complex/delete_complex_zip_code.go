package mysqlcomplexadapter

import (
	"context"
	"database/sql"
)

func (ca *ComplexAdapter) DeleteComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	query := "DELETE FROM complex_zip_codes WHERE id = ? LIMIT 1;"
	return ca.Delete(ctx, tx, query, id)
}
