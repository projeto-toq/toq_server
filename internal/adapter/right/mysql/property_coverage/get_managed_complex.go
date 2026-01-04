package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetManagedComplex fetches a complex by id and kind; returns sql.ErrNoRows when not found.
func (a *PropertyCoverageAdapter) GetManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	switch kind {
	case propertycoveragemodel.CoverageKindVertical:
		return a.getVerticalComplexByID(ctx, tx, id)
	case propertycoveragemodel.CoverageKindHorizontal:
		return a.getHorizontalComplexByID(ctx, tx, id)
	case propertycoveragemodel.CoverageKindStandalone:
		return a.getStandaloneComplexByID(ctx, tx, id)
	default:
		return nil, sql.ErrNoRows
	}
}

func (a *PropertyCoverageAdapter) getVerticalComplexByID(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
        SELECT vc.id, vc.name, vc.zip_code, vc.street, vc.number, vc.neighborhood, vc.city, vc.state,
               vc.reception_phone, vc.sector, vc.main_registration, vc.type
        FROM vertical_complexes vc
        WHERE vc.id = ?
        LIMIT 1;
    `

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{id}, propertycoveragemodel.CoverageKindVertical)
	if err != nil {
		return nil, err
	}

	domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)

	towers, err := a.listTowers(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetTowers(towers)

	sizes, err := a.listSizes(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetSizes(sizes)

	return domain, nil
}

func (a *PropertyCoverageAdapter) getHorizontalComplexByID(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
        SELECT hc.id, hc.name, hc.zip_code, hc.street, hc.number, hc.neighborhood, hc.city, hc.state,
               hc.reception_phone, hc.sector, hc.main_registration, hc.type
        FROM horizontal_complexes hc
        WHERE hc.id = ?
        LIMIT 1;
    `

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{id}, propertycoveragemodel.CoverageKindHorizontal)
	if err != nil {
		return nil, err
	}

	domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)

	zipCodes, err := a.listZipCodes(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetZipCodes(zipCodes)

	return domain, nil
}

func (a *PropertyCoverageAdapter) getStandaloneComplexByID(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
        SELECT nc.id, NULL AS name, nc.zip_code, NULL AS street, NULL AS number,
               nc.neighborhood, nc.city, nc.state, NULL AS reception_phone, nc.sector,
               NULL AS main_registration, nc.type
        FROM no_complex_zip_codes nc
        WHERE nc.id = ?
        LIMIT 1;
    `

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{id}, propertycoveragemodel.CoverageKindStandalone)
	if err != nil {
		return nil, err
	}

	return propertycoverageconverters.ManagedComplexEntityToDomain(entity), nil
}
