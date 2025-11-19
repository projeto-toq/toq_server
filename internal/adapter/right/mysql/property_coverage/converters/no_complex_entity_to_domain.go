package propertycoverageconverters

import (
	"fmt"

	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// NoComplexEntityToDomain converts the standalone coverage entity into the domain object.
func NoComplexEntityToDomain(entity propertycoverageentities.NoComplexCoverageEntity) (propertycoveragemodel.CoverageInterface, error) {
	coverage := propertycoveragemodel.NewCoverage()
	coverage.SetSource(propertycoveragemodel.CoverageSourceStandalone)

	if entity.PropertyTypesBitmask < 0 {
		return nil, fmt.Errorf("negative property type bitmask for standalone coverage: %d", entity.PropertyTypesBitmask)
	}

	coverage.SetPropertyTypes(globalmodel.PropertyType(entity.PropertyTypesBitmask))
	return coverage, nil
}
