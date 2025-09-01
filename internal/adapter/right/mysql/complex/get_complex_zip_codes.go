package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	complexrepoconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexZipCodes(ctx context.Context, tx *sql.Tx, complexID int64) (complexZipCodes []complexmodel.ComplexZipCodeInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM complex_zip_codes WHERE complex_id = ?;`

	entities, err := ca.Read(ctx, tx, query, complexID)
	if err != nil {
		slog.Error("mysqlcomplexadapter/GetComplexZipCodes: error executing Read", "error", err)
		return nil, fmt.Errorf("get complex zip codes read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, entity := range entities {
		complexZipCode, err := complexrepoconverters.ComplexZipCodeEntityToDomain(entity)
		if err != nil {
			return nil, fmt.Errorf("convert complex zip code entity: %w", err)
		}

		complexZipCodes = append(complexZipCodes, complexZipCode)
	}

	return
}
