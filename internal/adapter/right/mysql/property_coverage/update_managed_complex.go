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

// UpdateManagedComplex updates a managed complex according to its kind; returns sql.ErrNoRows when not found.
func (a *PropertyCoverageAdapter) UpdateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	writeEntity := propertycoverageconverters.ManagedComplexDomainToEntity(entity)

	switch entity.Kind() {
	case propertycoveragemodel.CoverageKindVertical:
		return a.updateVerticalComplex(ctx, tx, writeEntity)
	case propertycoveragemodel.CoverageKindHorizontal:
		return a.updateHorizontalComplex(ctx, tx, writeEntity)
	case propertycoveragemodel.CoverageKindStandalone:
		return a.updateStandaloneComplex(ctx, tx, writeEntity)
	default:
		return 0, fmt.Errorf("unsupported coverage kind %s", entity.Kind())
	}
}

func (a *PropertyCoverageAdapter) updateVerticalComplex(ctx context.Context, tx *sql.Tx, entity propertycoverageentities.ManagedComplexEntity) (int64, error) {
	const query = `
        UPDATE vertical_complexes SET
            name = ?,
            zip_code = ?,
            street = ?,
            number = ?,
            neighborhood = ?,
            city = ?,
            state = ?,
            reception_phone = ?,
            sector = ?,
            main_registration = ?,
            type = ?
        WHERE id = ?
        LIMIT 1;
    `

	return a.execUpdate(ctx, tx, "update", query,
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
		entity.ID,
	)
}

func (a *PropertyCoverageAdapter) updateHorizontalComplex(ctx context.Context, tx *sql.Tx, entity propertycoverageentities.ManagedComplexEntity) (int64, error) {
	const query = `
        UPDATE horizontal_complexes SET
            name = ?,
            zip_code = ?,
            street = ?,
            number = ?,
            neighborhood = ?,
            city = ?,
            state = ?,
            reception_phone = ?,
            sector = ?,
            main_registration = ?,
            type = ?
        WHERE id = ?
        LIMIT 1;
    `

	return a.execUpdate(ctx, tx, "update", query,
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
		entity.ID,
	)
}

func (a *PropertyCoverageAdapter) updateStandaloneComplex(ctx context.Context, tx *sql.Tx, entity propertycoverageentities.ManagedComplexEntity) (int64, error) {
	const query = `
        UPDATE no_complex_zip_codes SET
            zip_code = ?,
            neighborhood = ?,
            city = ?,
            state = ?,
            sector = ?,
            type = ?
        WHERE id = ?
        LIMIT 1;
    `

	return a.execUpdate(ctx, tx, "update", query,
		entity.ZipCode,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.Sector,
		entity.PropertyTypes,
		entity.ID,
	)
}
