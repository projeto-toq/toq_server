package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateHorizontalComplexZipCode updates a zip mapping; returns sql.ErrNoRows when no row is affected.
func (a *PropertyCoverageAdapter) UpdateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	entity := propertycoverageconverters.HorizontalComplexZipCodeDomainToEntity(zip)

	const query = `
        UPDATE horizontal_complex_zip_codes SET
            horizontal_complex_id = ?,
            zip_code = ?
        WHERE id = ?
        LIMIT 1;
    `

	return a.execUpdate(ctx, tx, "update", query, entity.HorizontalComplexID, entity.ZipCode, entity.ID)
}
