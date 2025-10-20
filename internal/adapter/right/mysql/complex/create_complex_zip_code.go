package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) CreateComplexZipCode(ctx context.Context, tx *sql.Tx, zip complexmodel.ComplexZipCodeInterface) (int64, error) {
	query := `INSERT INTO complex_zip_codes (
		complex_id,
		zip_code
	) VALUES (?, ?);`

	return ca.Create(
		ctx,
		tx,
		query,
		zip.ComplexID(),
		zip.ZipCode(),
	)
}
