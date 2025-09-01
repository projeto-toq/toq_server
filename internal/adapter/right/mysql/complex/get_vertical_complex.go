package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	complexrepoconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetVerticalComplex(ctx context.Context, tx *sql.Tx, zipCode string, number string) (complex complexmodel.ComplexInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	query := `SELECT * FROM complex WHERE zip_code = ? AND number = ?;`

	entity, err := ca.Read(ctx, tx, query, zipCode)
	if err != nil {
		slog.Error("mysqlcomplexadapter/GetVerticalComplex: error executing Read", "error", err)
		return nil, fmt.Errorf("get vertical complex read: %w", err)
	}

	if len(entity) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entity) > 1 {
		return nil, errors.New("multiple vertical complex rows found")
	}

	complex, err = complexrepoconverters.ComplexEntityToDomain(entity[0])
	if err != nil {
		return
	}

	complexSizes, err := ca.GetComplexSizes(ctx, tx, complex.ID())
	if err != nil {
		return
	}
	complex.SetComplexSizes(complexSizes)

	complexTowers, err := ca.GetComplexTowers(ctx, tx, complex.ID())
	if err != nil {
		return
	}
	complex.SetComplexTowers(complexTowers)

	complexZipCodes, err := ca.GetComplexZipCodes(ctx, tx, complex.ID())
	if err != nil {
		return
	}
	complex.SetComplexZipCodes(complexZipCodes)

	return
}
