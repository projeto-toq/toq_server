package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateHorizontalComplexZipCode inserts a zip code mapping for a horizontal complex and returns the created id.
func (a *PropertyCoverageAdapter) CreateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	entity := propertycoverageconverters.HorizontalComplexZipCodeDomainToEntity(zip)

	const query = `
        INSERT INTO horizontal_complex_zip_codes (
            horizontal_complex_id, zip_code
        ) VALUES (?, ?);
    `

	return a.execInsert(ctx, tx, query, entity.HorizontalComplexID, entity.ZipCode)
}
