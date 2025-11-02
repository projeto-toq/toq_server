package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) CreateComplexSize(ctx context.Context, tx *sql.Tx, size complexmodel.ComplexSizeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO complex_sizes (
		complex_id,
		size,
		description
	) VALUES (?, ?, ?);`

	result, err := ca.ExecContext(
		ctx,
		tx,
		"insert",
		query,
		size.ComplexID(),
		size.Size(),
		nullableStringValue(size.Description()),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.size.create.exec_error", "error", err, "complex_id", size.ComplexID())
		return 0, fmt.Errorf("insert complex size: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.size.create.last_insert_id_error", "error", err, "complex_id", size.ComplexID())
		return 0, fmt.Errorf("complex size last insert id: %w", err)
	}

	return id, nil
}
