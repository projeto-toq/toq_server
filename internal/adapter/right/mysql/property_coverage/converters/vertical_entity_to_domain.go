package propertycoverageconverters

import (
	"fmt"

	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// VerticalEntityToDomain converts a DB entity into a domain coverage object.
func VerticalEntityToDomain(entity propertycoverageentities.VerticalCoverageEntity) (propertycoveragemodel.CoverageInterface, error) {
	coverage := propertycoveragemodel.NewCoverage()
	coverage.SetSource(propertycoveragemodel.CoverageSourceVertical)
	coverage.SetComplexName(entity.Name)
	coverage.SetMainRegistration(entity.MainRegistration)

	if entity.PropertyTypesBitmask < 0 {
		return nil, fmt.Errorf("negative property type bitmask for vertical complex: %d", entity.PropertyTypesBitmask)
	}

	coverage.SetPropertyTypes(globalmodel.PropertyType(entity.PropertyTypesBitmask))
	return coverage, nil
}
