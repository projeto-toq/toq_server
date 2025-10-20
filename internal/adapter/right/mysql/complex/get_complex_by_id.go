package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, name, zip_code, street, number, neighborhood, city, state, reception_phone, sector, main_registration, type
	FROM complex WHERE id = ? LIMIT 1;`

	entities, err := ca.Read(ctx, tx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_by_id.read_error", "error", err, "id", id)
		return nil, fmt.Errorf("get complex by id read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	complexDomain, err := complexrepoconverters.ComplexEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_by_id.convert_error", "error", err, "id", id)
		return nil, fmt.Errorf("convert complex entity: %w", err)
	}

	sizes, err := ca.GetComplexSizes(ctx, tx, complexDomain.ID())
	if err != nil && err != sql.ErrNoRows {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_by_id.sizes_error", "error", err, "complex_id", complexDomain.ID())
		return nil, fmt.Errorf("load complex sizes: %w", err)
	}
	if err == nil {
		complexDomain.SetComplexSizes(sizes)
	}

	towers, err := ca.GetComplexTowers(ctx, tx, complexDomain.ID())
	if err != nil && err != sql.ErrNoRows {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_by_id.towers_error", "error", err, "complex_id", complexDomain.ID())
		return nil, fmt.Errorf("load complex towers: %w", err)
	}
	if err == nil {
		complexDomain.SetComplexTowers(towers)
	}

	zipCodes, err := ca.GetComplexZipCodes(ctx, tx, complexDomain.ID())
	if err != nil && err != sql.ErrNoRows {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_by_id.zip_codes_error", "error", err, "complex_id", complexDomain.ID())
		return nil, fmt.Errorf("load complex zip codes: %w", err)
	}
	if err == nil {
		complexDomain.SetComplexZipCodes(zipCodes)
	}

	return complexDomain, nil
}
