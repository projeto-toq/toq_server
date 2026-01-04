package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetVerticalComplexByZipNumber busca complexo vertical por CEP + numero e retorna torres/tamanhos; sql.ErrNoRows se não encontrar.
func (a *PropertyCoverageAdapter) GetVerticalComplexByZipNumber(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	const query = `
        SELECT vc.id, vc.name, vc.zip_code, vc.street, vc.number, vc.neighborhood, vc.city, vc.state,
               vc.reception_phone, vc.sector, vc.main_registration, vc.type
        FROM vertical_complexes vc
        WHERE vc.zip_code = ?
          AND (
                UPPER(REPLACE(TRIM(vc.number), ' ', '')) = ?
             OR FIND_IN_SET(
                  ?,
                  REPLACE(REPLACE(UPPER(TRIM(vc.number)), ' ', ''), ';', ',')
                ) > 0
              )
        LIMIT 1;
    `

	entity, fetchErr := a.fetchManagedComplex(ctx, tx, query, []any{zipCode, number, number}, propertycoveragemodel.CoverageKindVertical)
	if fetchErr != nil {
		return nil, fetchErr
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
