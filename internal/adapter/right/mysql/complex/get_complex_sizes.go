package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexrepoconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexSizes(ctx context.Context, tx *sql.Tx, complexID int64) (complexSizes []complexmodel.ComplexSizeInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM complex_sizes WHERE complex_id = ?;`

	entities, err := ca.Read(ctx, tx, query, complexID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_sizes.read_error", "error", err, "complex_id", complexID)
		return nil, fmt.Errorf("get complex sizes read: %w", err)
	}

	if len(entities) == 0 {
		err = sql.ErrNoRows
		utils.SetSpanError(ctx, err)
		logger.Warn("mysql.complex.get_sizes.not_found", "complex_id", complexID)
		return nil, err
	}

	for _, entity := range entities {
		complexSize, errConv := complexrepoconverters.ComplexSizeEntityToDomain(entity)
		if errConv != nil {
			utils.SetSpanError(ctx, errConv)
			logger.Error("mysql.complex.get_sizes.convert_error", "error", errConv, "complex_id", complexID)
			return nil, fmt.Errorf("convert complex size entity: %w", errConv)
		}

		complexSizes = append(complexSizes, complexSize)
	}

	return
}
