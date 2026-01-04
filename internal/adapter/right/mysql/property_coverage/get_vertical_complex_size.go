package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetVerticalComplexSize fetches a size by id; returns sql.ErrNoRows when not found.
func (a *PropertyCoverageAdapter) GetVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	const query = `
        SELECT id, vertical_complex_id, size, description
        FROM vertical_complex_sizes
        WHERE id = ?
        LIMIT 1;
    `

	row := a.QueryRowContext(ctx, tx, "select", query, id)
	var entity propertycoverageentities.VerticalComplexSizeEntity
	if scanErr := row.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Size, &entity.Description); scanErr != nil {
		if scanErr != sql.ErrNoRows {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.property_coverage.size.get.scan_error", "err", scanErr)
		}
		return nil, scanErr
	}

	return propertycoverageconverters.VerticalComplexSizeEntityToDomain(entity), nil
}
