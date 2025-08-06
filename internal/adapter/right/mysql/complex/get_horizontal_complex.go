package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"log/slog"

	complexrepoconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ca *ComplexAdapter) GetHorizontalComplex(ctx context.Context, tx *sql.Tx, zipCode string) (complex complexmodel.ComplexInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT c.*
				FROM complex c
				JOIN complex_zip_codes z on z.complex_id = c.id
				WHERE z.zip_code = ?;`

	entities, err := ca.Read(ctx, tx, query, zipCode)
	if err != nil {
		slog.Error("mysqlcomplexadapter/GetHorizontalComplex: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "No complex found for the provided zip code")
	}

	for _, entity := range entities {
		complex, err = complexrepoconverters.ComplexEntityToDomain(entity)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal server error")
		}
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
