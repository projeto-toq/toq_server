package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetHorizontalComplex(ctx context.Context, tx *sql.Tx, zipCode string) (complex complexmodel.ComplexInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT c.*
				FROM complex c
				JOIN complex_zip_codes z on z.complex_id = c.id
				WHERE z.zip_code = ?;`

	entities, err := ca.Read(ctx, tx, query, zipCode)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_horizontal.read_error", "error", err, "zip_code", zipCode)
		return nil, fmt.Errorf("get horizontal complex read: %w", err)
	}

	if len(entities) == 0 {
		err = sql.ErrNoRows
		utils.SetSpanError(ctx, err)
		logger.Warn("mysql.complex.get_horizontal.not_found", "zip_code", zipCode)
		return nil, err
	}

	for _, entity := range entities {
		complex, err = complexrepoconverters.ComplexEntityToDomain(entity)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.complex.get_horizontal.convert_error", "error", err, "zip_code", zipCode)
			return nil, fmt.Errorf("convert complex entity: %w", err)
		}
	}

	// complexSizes, err := ca.GetComplexSizes(ctx, tx, complex.ID())
	// if err != nil {
	// 	utils.SetSpanError(ctx, err)
	// 	logger.Error("mysql.complex.get_horizontal.sizes_error", "error", err, "complex_id", complex.ID())
	// 	return
	// }
	// complex.SetComplexSizes(complexSizes)

	// complexTowers, err := ca.GetComplexTowers(ctx, tx, complex.ID())
	// if err != nil {
	// 	utils.SetSpanError(ctx, err)
	// 	logger.Error("mysql.complex.get_horizontal.towers_error", "error", err, "complex_id", complex.ID())
	// 	return
	// }
	// complex.SetComplexTowers(complexTowers)

	complexZipCodes, err := ca.GetComplexZipCodes(ctx, tx, complex.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_horizontal.zip_codes_error", "error", err, "complex_id", complex.ID())
		return
	}
	complex.SetComplexZipCodes(complexZipCodes)

	return
}
