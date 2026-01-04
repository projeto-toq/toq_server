package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"
	"fmt"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateManagedComplex inserts a complex according to the informed CoverageKind and returns the created id.
func (a *PropertyCoverageAdapter) CreateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	writeEntity := propertycoverageconverters.ManagedComplexDomainToEntity(entity)

	switch entity.Kind() {
	case propertycoveragemodel.CoverageKindVertical:
		return a.insertVerticalComplex(ctx, tx, writeEntity)
	case propertycoveragemodel.CoverageKindHorizontal:
		return a.insertHorizontalComplex(ctx, tx, writeEntity)
	case propertycoveragemodel.CoverageKindStandalone:
		return a.insertStandaloneComplex(ctx, tx, writeEntity)
	default:
		return 0, fmt.Errorf("unsupported coverage kind %s", entity.Kind())
	}
}

func (a *PropertyCoverageAdapter) insertVerticalComplex(ctx context.Context, tx *sql.Tx, entity propertycoverageentities.ManagedComplexEntity) (int64, error) {
	const query = `
        INSERT INTO vertical_complexes (
            name, zip_code, street, number, neighborhood, city, state,
            reception_phone, sector, main_registration, type
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	return a.execInsert(ctx, tx, query,
		entity.Name,
		entity.ZipCode,
		entity.Street,
		entity.Number,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.ReceptionPhone,
		entity.Sector,
		entity.MainRegistration,
		entity.PropertyTypes,
	)
}

func (a *PropertyCoverageAdapter) insertHorizontalComplex(ctx context.Context, tx *sql.Tx, entity propertycoverageentities.ManagedComplexEntity) (int64, error) {
	const query = `
        INSERT INTO horizontal_complexes (
            name, zip_code, street, number, neighborhood, city, state,
            reception_phone, sector, main_registration, type
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	return a.execInsert(ctx, tx, query,
		entity.Name,
		entity.ZipCode,
		entity.Street,
		entity.Number,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.ReceptionPhone,
		entity.Sector,
		entity.MainRegistration,
		entity.PropertyTypes,
	)
}

func (a *PropertyCoverageAdapter) insertStandaloneComplex(ctx context.Context, tx *sql.Tx, entity propertycoverageentities.ManagedComplexEntity) (int64, error) {
	const query = `
        INSERT INTO no_complex_zip_codes (
            zip_code, neighborhood, city, state, sector, type
        ) VALUES (?, ?, ?, ?, ?, ?);
    `

	return a.execInsert(ctx, tx, query,
		entity.ZipCode,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.Sector,
		entity.PropertyTypes,
	)
}
