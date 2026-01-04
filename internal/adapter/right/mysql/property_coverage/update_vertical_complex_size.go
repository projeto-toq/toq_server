package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateVerticalComplexSize updates a size; returns sql.ErrNoRows when no row is affected.
func (a *PropertyCoverageAdapter) UpdateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	entity := propertycoverageconverters.VerticalComplexSizeDomainToEntity(size)

	const query = `
        UPDATE vertical_complex_sizes SET
            vertical_complex_id = ?,
            size = ?,
            description = ?
        WHERE id = ?
        LIMIT 1;
    `

	return a.execUpdate(ctx, tx, "update", query,
		entity.VerticalComplexID,
		entity.Size,
		entity.Description,
		entity.ID,
	)
}
