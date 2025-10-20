package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	repository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) ListComplexZipCodes(ctx context.Context, tx *sql.Tx, params repository.ListComplexZipCodesParams) ([]complexmodel.ComplexZipCodeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	var builder strings.Builder
	builder.WriteString("SELECT id, complex_id, zip_code FROM complex_zip_codes WHERE 1=1")
	args := make([]any, 0)

	if params.ComplexID > 0 {
		builder.WriteString(" AND complex_id = ?")
		args = append(args, params.ComplexID)
	}

	if params.ZipCode != "" {
		builder.WriteString(" AND zip_code = ?")
		args = append(args, params.ZipCode)
	}

	builder.WriteString(" ORDER BY id ASC")

	if params.Limit > 0 {
		builder.WriteString(" LIMIT ?")
		args = append(args, params.Limit)
	}

	if params.Offset > 0 {
		builder.WriteString(" OFFSET ?")
		args = append(args, params.Offset)
	}

	query := builder.String()

	entities, err := ca.Read(ctx, tx, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.list_zip_codes.read_error", "error", err, "params", params)
		return nil, fmt.Errorf("list complex zip codes read: %w", err)
	}

	zipCodes := make([]complexmodel.ComplexZipCodeInterface, 0, len(entities))

	for _, entity := range entities {
		zipCode, errConv := complexrepoconverters.ComplexZipCodeEntityToDomain(entity)
		if errConv != nil {
			utils.SetSpanError(ctx, errConv)
			logger.Error("mysql.complex.list_zip_codes.convert_error", "error", errConv)
			return nil, fmt.Errorf("convert complex zip code entity: %w", errConv)
		}

		zipCodes = append(zipCodes, zipCode)
	}

	return zipCodes, nil
}
