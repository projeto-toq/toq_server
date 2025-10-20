package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) UpdateComplexSize(ctx context.Context, tx *sql.Tx, size complexmodel.ComplexSizeInterface) (int64, error) {
	query := `UPDATE complex_sizes SET
		size = ?,
		description = ?
	WHERE id = ?
	AND complex_id = ?
	LIMIT 1;`

	return ca.Update(
		ctx,
		tx,
		query,
		size.Size(),
		nullableStringValue(size.Description()),
		size.ID(),
		size.ComplexID(),
	)
}
