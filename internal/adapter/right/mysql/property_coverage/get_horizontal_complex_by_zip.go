package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetHorizontalComplexByZip retrieves a horizontal complex by its zip code.
// It searches in both the main complex record and the associated zip codes table.
// Returns sql.ErrNoRows if no complex is found associated with the provided zip code.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction
//   - zipCode: The zip code to search for (numbers only)
//
// Returns:
//   - complex: The managed complex domain object
//   - error: sql.ErrNoRows or infrastructure error
func (a *PropertyCoverageAdapter) GetHorizontalComplexByZip(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.ManagedComplexInterface, error) {
	// Initialize tracing as per Section 7.3
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)

	// Query checks both main zip_code and associated zip_codes using LEFT JOIN
	// DISTINCT is used to avoid duplicates if zip matches both (unlikely but safe)
	const query = `
		SELECT DISTINCT hc.id, hc.name, hc.zip_code, hc.street, hc.number, hc.neighborhood, hc.city, hc.state,
		       hc.reception_phone, hc.sector, hc.main_registration, hc.type
		FROM horizontal_complexes hc
		LEFT JOIN horizontal_complex_zip_codes hcz ON hcz.horizontal_complex_id = hc.id
		WHERE hc.zip_code = ? OR hcz.zip_code = ?
		LIMIT 1;
	`

	// Reuse existing private helper fetchManagedComplex from the same package
	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{zipCode, zipCode}, propertycoveragemodel.CoverageKindHorizontal)
	if err != nil {
		// fetchManagedComplex handles logging and SetSpanError
		return nil, err
	}

	// Convert entity to domain
	domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)

	// Populate associated zip codes using existing helper
	zipCodes, err := a.listZipCodes(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetZipCodes(zipCodes)

	return domain, nil
}
