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

func (ca *ComplexAdapter) GetComplexSizes(ctx context.Context, tx *sql.Tx, complexID int64) (complexSizes []complexmodel.ComplexSizeInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM complex_sizes WHERE complex_id = ?;`

	entities, err := ca.Read(ctx, tx, query, complexID)
	if err != nil {
		slog.Error("mysqlcomplexadapter/GetComplexSizes: error executing Read", "error", err)
		return nil, fmt.Errorf("get complex sizes read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, entity := range entities {
		complexSize, err1 := complexrepoconverters.ComplexSizeEntityToDomain(entity)
		if err1 != nil {
			return nil, fmt.Errorf("convert complex size entity: %w", err1)
		}

		complexSizes = append(complexSizes, complexSize)
	}

	return
}
