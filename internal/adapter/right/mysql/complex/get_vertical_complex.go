package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetVerticalComplex(ctx context.Context, tx *sql.Tx, zipCode string, number string) (complex complexmodel.ComplexInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM complex WHERE zip_code = ? AND number = ?;`

	entity, err := ca.Read(ctx, tx, query, zipCode)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_vertical.read_error", "error", err, "zip_code", zipCode, "number", number)
		return nil, fmt.Errorf("get vertical complex read: %w", err)
	}

	if len(entity) == 0 {
		err = sql.ErrNoRows
		utils.SetSpanError(ctx, err)
		logger.Warn("mysql.complex.get_vertical.not_found", "zip_code", zipCode, "number", number)
		return nil, err
	}

	if len(entity) > 1 {
		err = errors.New("multiple vertical complex rows found")
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_vertical.multiple_rows", "error", err, "zip_code", zipCode, "number", number)
		return nil, err
	}

	complex, err = complexrepoconverters.ComplexEntityToDomain(entity[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_vertical.convert_error", "error", err, "zip_code", zipCode, "number", number)
		return
	}

	complexSizes, err := ca.GetComplexSizes(ctx, tx, complex.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_vertical.sizes_error", "error", err, "complex_id", complex.ID())
		return
	}
	complex.SetComplexSizes(complexSizes)

	complexTowers, err := ca.GetComplexTowers(ctx, tx, complex.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_vertical.towers_error", "error", err, "complex_id", complex.ID())
		return
	}
	complex.SetComplexTowers(complexTowers)

	complexZipCodes, err := ca.GetComplexZipCodes(ctx, tx, complex.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_vertical.zip_codes_error", "error", err, "complex_id", complex.ID())
		return
	}
	complex.SetComplexZipCodes(complexZipCodes)

	return
}
