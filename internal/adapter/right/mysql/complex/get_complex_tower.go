package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"log/slog"

	complexrepoconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexTowers(ctx context.Context, tx *sql.Tx, complexID int64) (complexTowers []complexmodel.ComplexTowerInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM complex_towers WHERE complex_id = ?;`

	entities, err := ca.Read(ctx, tx, query, complexID)
	if err != nil {
		slog.Error("mysqlcomplexadapter/GetComplexTowers: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	for _, entity := range entities {
		complexTower, err := complexrepoconverters.ComplexTowerEntityToDomain(entity)
		if err != nil {
			return nil, utils.ErrInternalServer
		}

		complexTowers = append(complexTowers, complexTower)
	}

	return
}
