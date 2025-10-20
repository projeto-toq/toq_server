package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) CreateComplexSize(ctx context.Context, tx *sql.Tx, size complexmodel.ComplexSizeInterface) (int64, error) {
	query := `INSERT INTO complex_sizes (
		complex_id,
		size,
		description
	) VALUES (?, ?, ?);`

	return ca.Create(
		ctx,
		tx,
		query,
		size.ComplexID(),
		size.Size(),
		nullableStringValue(size.Description()),
	)
}
