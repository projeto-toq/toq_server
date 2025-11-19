package propertycoverageconverters

import (
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

func ManagedComplexEntityToDomain(entity propertycoverageentities.ManagedComplexEntity) propertycoveragemodel.ManagedComplexInterface {
	domain := propertycoveragemodel.NewManagedComplex()
	domain.SetID(entity.ID)
	domain.SetKind(entity.Kind)
	domain.SetName(entity.Name.String)
	domain.SetZipCode(entity.ZipCode)
	domain.SetStreet(entity.Street.String)
	domain.SetNumber(entity.Number.String)
	domain.SetNeighborhood(entity.Neighborhood.String)
	domain.SetCity(entity.City)
	domain.SetState(entity.State)
	domain.SetReceptionPhone(entity.ReceptionPhone.String)
	domain.SetSector(propertycoveragemodel.Sector(entity.Sector))
	domain.SetMainRegistration(entity.MainRegistration.String)
	domain.SetPropertyTypes(globalmodel.PropertyType(entity.PropertyTypes))
	return domain
}

func VerticalComplexTowerEntityToDomain(entity propertycoverageentities.VerticalComplexTowerEntity) propertycoveragemodel.VerticalComplexTowerInterface {
	tower := propertycoveragemodel.NewVerticalComplexTower()
	tower.SetID(entity.ID)
	tower.SetVerticalComplexID(entity.VerticalComplexID)
	tower.SetTower(entity.Tower)
	if ptr := intToOptional(entity.Floors); ptr != nil {
		tower.SetFloors(ptr)
	}
	if ptr := intToOptional(entity.TotalUnits); ptr != nil {
		tower.SetTotalUnits(ptr)
	}
	if ptr := intToOptional(entity.UnitsPerFloor); ptr != nil {
		tower.SetUnitsPerFloor(ptr)
	}
	return tower
}

func VerticalComplexSizeEntityToDomain(entity propertycoverageentities.VerticalComplexSizeEntity) propertycoveragemodel.VerticalComplexSizeInterface {
	size := propertycoveragemodel.NewVerticalComplexSize()
	size.SetID(entity.ID)
	size.SetVerticalComplexID(entity.VerticalComplexID)
	size.SetSize(entity.Size)
	size.SetDescription(entity.Description.String)
	return size
}

func HorizontalComplexZipCodeEntityToDomain(entity propertycoverageentities.HorizontalComplexZipCodeEntity) propertycoveragemodel.HorizontalComplexZipCodeInterface {
	zip := propertycoveragemodel.NewHorizontalComplexZipCode()
	zip.SetID(entity.ID)
	zip.SetHorizontalComplexID(entity.HorizontalComplexID)
	zip.SetZipCode(entity.ZipCode)
	return zip
}

func intToOptional(value int) *int {
	if value == 0 {
		return nil
	}
	copy := value
	return &copy
}
