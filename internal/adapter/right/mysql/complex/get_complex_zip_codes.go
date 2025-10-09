package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexZipCodes(ctx context.Context, tx *sql.Tx, complexID int64) (complexZipCodes []complexmodel.ComplexZipCodeInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM complex_zip_codes WHERE complex_id = ?;`

	entities, err := ca.Read(ctx, tx, query, complexID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_zip_codes.read_error", "error", err, "complex_id", complexID)
		return nil, fmt.Errorf("get complex zip codes read: %w", err)
	}

	if len(entities) == 0 {
		err = sql.ErrNoRows
		utils.SetSpanError(ctx, err)
		logger.Warn("mysql.complex.get_zip_codes.not_found", "complex_id", complexID)
		return nil, err
	}

	for _, entity := range entities {
		complexZipCode, errConv := complexrepoconverters.ComplexZipCodeEntityToDomain(entity)
		if errConv != nil {
			utils.SetSpanError(ctx, errConv)
			logger.Error("mysql.complex.get_zip_codes.convert_error", "error", errConv, "complex_id", complexID)
			return nil, fmt.Errorf("convert complex zip code entity: %w", errConv)
		}

		complexZipCodes = append(complexZipCodes, complexZipCode)
	}

	return
}
