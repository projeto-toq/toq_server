package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) UpdateComplexZipCode(ctx context.Context, tx *sql.Tx, zip complexmodel.ComplexZipCodeInterface) (int64, error) {
	query := `UPDATE complex_zip_codes SET
		zip_code = ?
	WHERE id = ?
	AND complex_id = ?
	LIMIT 1;`

	return ca.Update(
		ctx,
		tx,
		query,
		zip.ZipCode(),
		zip.ID(),
		zip.ComplexID(),
	)
}
