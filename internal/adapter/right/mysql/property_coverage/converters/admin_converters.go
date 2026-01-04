package propertycoverageconverters

import (
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// ManagedComplexEntityToDomain converte entidade admin para domínio ManagedComplexInterface.
// Campos sql.Null* são checados antes de preencher o domínio para evitar vazios incorretos.
func ManagedComplexEntityToDomain(entity propertycoverageentities.ManagedComplexEntity) propertycoveragemodel.ManagedComplexInterface {
	domain := propertycoveragemodel.NewManagedComplex()
	domain.SetID(entity.ID)
	domain.SetKind(entity.Kind)
	if entity.Name.Valid {
		domain.SetName(entity.Name.String)
	}
	domain.SetZipCode(entity.ZipCode)
	if entity.Street.Valid {
		domain.SetStreet(entity.Street.String)
	}
	if entity.Number.Valid {
		domain.SetNumber(entity.Number.String)
	}
	if entity.Neighborhood.Valid {
		domain.SetNeighborhood(entity.Neighborhood.String)
	}
	domain.SetCity(entity.City)
	domain.SetState(entity.State)
	if entity.ReceptionPhone.Valid {
		domain.SetReceptionPhone(entity.ReceptionPhone.String)
	}
	domain.SetSector(propertycoveragemodel.Sector(entity.Sector))
	if entity.MainRegistration.Valid {
		domain.SetMainRegistration(entity.MainRegistration.String)
	}
	domain.SetPropertyTypes(globalmodel.PropertyType(entity.PropertyTypes))
	return domain
}

// VerticalComplexTowerEntityToDomain converte entidade de torre para domínio.
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

// VerticalComplexSizeEntityToDomain converte entidade de tamanho para domínio.
func VerticalComplexSizeEntityToDomain(entity propertycoverageentities.VerticalComplexSizeEntity) propertycoveragemodel.VerticalComplexSizeInterface {
	size := propertycoveragemodel.NewVerticalComplexSize()
	size.SetID(entity.ID)
	size.SetVerticalComplexID(entity.VerticalComplexID)
	size.SetSize(entity.Size)
	if entity.Description.Valid {
		size.SetDescription(entity.Description.String)
	}
	return size
}

// HorizontalComplexZipCodeEntityToDomain converte entidade de CEP horizontal para domínio.
func HorizontalComplexZipCodeEntityToDomain(entity propertycoverageentities.HorizontalComplexZipCodeEntity) propertycoveragemodel.HorizontalComplexZipCodeInterface {
	zip := propertycoveragemodel.NewHorizontalComplexZipCode()
	zip.SetID(entity.ID)
	zip.SetHorizontalComplexID(entity.HorizontalComplexID)
	zip.SetZipCode(entity.ZipCode)
	return zip
}

// intToOptional transforma 0 em nil para campos numéricos opcionais que usam default 0 no banco.
func intToOptional(value int) *int {
	if value == 0 {
		return nil
	}
	copy := value
	return &copy
}
