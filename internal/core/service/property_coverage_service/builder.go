package propertycoverageservice

import (
	"fmt"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func buildManagedComplexFromCreateInput(input CreateComplexInput) (propertycoveragemodel.ManagedComplexInterface, error) {
	if err := validateCoverageKind(input.Kind); err != nil {
		return nil, err
	}

	if err := validateSector(input.Sector); err != nil {
		return nil, err
	}

	if input.PropertyType <= 0 {
		return nil, utils.ValidationError("propertyType", "Property type must be provided.")
	}

	normalizedZip, err := normalizeAndValidateZip(input.ZipCode)
	if err != nil {
		return nil, err
	}

	city := sanitizeString(input.City)
	if city == "" {
		return nil, utils.ValidationError("city", "City is required.")
	}

	state := sanitizeString(input.State)
	if state == "" {
		return nil, utils.ValidationError("state", "State is required.")
	}

	street := sanitizeString(input.Street)
	switch input.Kind {
	case propertycoveragemodel.CoverageKindVertical:
		if err := ensureCreateNameAndNumber(input); err != nil {
			return nil, err
		}
		if street == "" {
			return nil, utils.ValidationError("street", "Street is required for vertical complexes.")
		}
	case propertycoveragemodel.CoverageKindHorizontal:
		if err := validateRequiredField("name", input.Name); err != nil {
			return nil, err
		}
		if street == "" {
			return nil, utils.ValidationError("street", "Street is required for horizontal complexes.")
		}
	case propertycoveragemodel.CoverageKindStandalone:
		// Standalone coverage does not persist name/street information.
	default:
		return nil, utils.ValidationError("coverageType", fmt.Sprintf("Unsupported coverage type %s", input.Kind))
	}

	domain := propertycoveragemodel.NewManagedComplex()
	domain.SetKind(input.Kind)
	domain.SetZipCode(normalizedZip)
	domain.SetName(sanitizeString(input.Name))
	domain.SetStreet(street)
	domain.SetNumber(sanitizeString(input.Number))
	domain.SetNeighborhood(normalizeOptional(input.Neighborhood))
	domain.SetCity(city)
	domain.SetState(state)
	domain.SetReceptionPhone(normalizeOptional(input.ReceptionPhone))
	domain.SetSector(input.Sector)
	domain.SetMainRegistration(normalizeOptional(input.MainRegistration))
	domain.SetPropertyTypes(input.PropertyType)

	return domain, nil
}

func ensureCreateNameAndNumber(input CreateComplexInput) error {
	if err := validateRequiredField("name", input.Name); err != nil {
		return err
	}
	if err := validateRequiredField("number", input.Number); err != nil {
		return err
	}
	return nil
}
