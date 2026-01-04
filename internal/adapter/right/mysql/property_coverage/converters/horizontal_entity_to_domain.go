package propertycoverageconverters

import (
	"fmt"

	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// HorizontalEntityToDomain converte entity de horizontal_complexes para domínio CoverageInterface.
// Trata main_registration nullable e valida bitmask (não pode ser negativo, pois schema é UNSIGNED).
func HorizontalEntityToDomain(entity propertycoverageentities.HorizontalCoverageEntity) (propertycoveragemodel.CoverageInterface, error) {
	coverage := propertycoveragemodel.NewCoverage()
	coverage.SetSource(propertycoveragemodel.CoverageSourceHorizontal)
	coverage.SetComplexName(entity.Name)
	if entity.MainRegistration.Valid {
		coverage.SetMainRegistration(entity.MainRegistration.String)
	}

	if entity.PropertyTypesBitmask < 0 {
		return nil, fmt.Errorf("negative property type bitmask for horizontal complex: %d", entity.PropertyTypesBitmask)
	}

	coverage.SetPropertyTypes(globalmodel.PropertyType(entity.PropertyTypesBitmask))
	return coverage, nil
}
