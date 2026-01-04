package propertycoverageconverters

import (
	"database/sql"
	"strings"

	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// ManagedComplexDomainToEntity maps the domain aggregate to a DB entity for write operations.
// Optional string fields become sql.NullString; property types and sector are kept as-is.
func ManagedComplexDomainToEntity(domain propertycoveragemodel.ManagedComplexInterface) propertycoverageentities.ManagedComplexEntity {
	return propertycoverageentities.ManagedComplexEntity{
		ID:               domain.ID(),
		Kind:             domain.Kind(),
		Name:             toNullString(domain.Name()),
		ZipCode:          strings.TrimSpace(domain.ZipCode()),
		Street:           toNullString(domain.Street()),
		Number:           toNullString(domain.Number()),
		Neighborhood:     toNullString(domain.Neighborhood()),
		City:             strings.TrimSpace(domain.City()),
		State:            strings.TrimSpace(domain.State()),
		ReceptionPhone:   toNullString(domain.ReceptionPhone()),
		Sector:           uint8(domain.Sector()),
		MainRegistration: toNullString(domain.MainRegistration()),
		PropertyTypes:    uint16(domain.PropertyTypes()),
	}
}

// VerticalComplexTowerDomainToEntity maps tower domain data to its DB entity.
func VerticalComplexTowerDomainToEntity(domain propertycoveragemodel.VerticalComplexTowerInterface) propertycoverageentities.VerticalComplexTowerEntity {
	return propertycoverageentities.VerticalComplexTowerEntity{
		ID:                domain.ID(),
		VerticalComplexID: domain.VerticalComplexID(),
		Tower:             strings.TrimSpace(domain.Tower()),
		Floors:            intOrZero(domain.Floors()),
		TotalUnits:        intOrZero(domain.TotalUnits()),
		UnitsPerFloor:     intOrZero(domain.UnitsPerFloor()),
	}
}

// VerticalComplexSizeDomainToEntity maps size domain data to its DB entity.
func VerticalComplexSizeDomainToEntity(domain propertycoveragemodel.VerticalComplexSizeInterface) propertycoverageentities.VerticalComplexSizeEntity {
	return propertycoverageentities.VerticalComplexSizeEntity{
		ID:                domain.ID(),
		VerticalComplexID: domain.VerticalComplexID(),
		Size:              domain.Size(),
		Description:       toNullString(domain.Description()),
	}
}

// HorizontalComplexZipCodeDomainToEntity maps zip domain data to its DB entity.
func HorizontalComplexZipCodeDomainToEntity(domain propertycoveragemodel.HorizontalComplexZipCodeInterface) propertycoverageentities.HorizontalComplexZipCodeEntity {
	return propertycoverageentities.HorizontalComplexZipCodeEntity{
		ID:                  domain.ID(),
		HorizontalComplexID: domain.HorizontalComplexID(),
		ZipCode:             strings.TrimSpace(domain.ZipCode()),
	}
}

func toNullString(value string) sql.NullString {
	trimmed := strings.TrimSpace(value)
	return sql.NullString{String: trimmed, Valid: trimmed != ""}
}

func intOrZero(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
