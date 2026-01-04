package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateVerticalComplexSize inserts a size for a vertical complex and returns the created id.
func (a *PropertyCoverageAdapter) CreateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	entity := propertycoverageconverters.VerticalComplexSizeDomainToEntity(size)

	const query = `
        INSERT INTO vertical_complex_sizes (
            vertical_complex_id, size, description
        ) VALUES (?, ?, ?);
    `

	return a.execInsert(ctx, tx, query,
		entity.VerticalComplexID,
		entity.Size,
		entity.Description,
	)
}
