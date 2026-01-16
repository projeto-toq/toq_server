package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// ToComplexResponse converte um empreendimento da cobertura para DTO HTTP.
func ToComplexResponse(entity propertycoveragemodel.ManagedComplexInterface) dto.ComplexResponse {
	resp := dto.ComplexResponse{
		CoverageType:     string(entity.Kind()),
		ID:               entity.ID(),
		Name:             entity.Name(),
		ZipCode:          entity.ZipCode(),
		Street:           entity.Street(),
		Number:           entity.Number(),
		Neighborhood:     entity.Neighborhood(),
		City:             entity.City(),
		State:            entity.State(),
		PhoneNumber:      entity.ReceptionPhone(),
		Sector:           uint8(entity.Sector()),
		MainRegistration: entity.MainRegistration(),
		PropertyType:     uint16(entity.PropertyTypes()),
	}

	if sizes := entity.Sizes(); len(sizes) > 0 {
		resp.Sizes = make([]dto.ComplexSizeResponse, 0, len(sizes))
		for _, size := range sizes {
			resp.Sizes = append(resp.Sizes, ToComplexSizeResponse(size))
		}
	}

	if towers := entity.Towers(); len(towers) > 0 {
		resp.Towers = make([]dto.ComplexTowerResponse, 0, len(towers))
		for _, tower := range towers {
			resp.Towers = append(resp.Towers, ToComplexTowerResponse(tower))
		}
	}

	if zipCodes := entity.ZipCodes(); len(zipCodes) > 0 {
		resp.ZipCodes = make([]dto.ComplexZipCodeResponse, 0, len(zipCodes))
		for _, zip := range zipCodes {
			resp.ZipCodes = append(resp.ZipCodes, ToComplexZipCodeResponse(zip))
		}
	}

	return resp
}

// ToComplexResponses converte vários empreendimentos para DTO.
func ToComplexResponses(entities []propertycoveragemodel.ManagedComplexInterface) []dto.ComplexResponse {
	responses := make([]dto.ComplexResponse, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, ToComplexResponse(entity))
	}
	return responses
}

// ToComplexTowerResponse converte uma torre do domínio em DTO.
func ToComplexTowerResponse(tower propertycoveragemodel.VerticalComplexTowerInterface) dto.ComplexTowerResponse {
	return dto.ComplexTowerResponse{
		ID:            tower.ID(),
		ComplexID:     tower.VerticalComplexID(),
		Tower:         tower.Tower(),
		Floors:        tower.Floors(),
		TotalUnits:    tower.TotalUnits(),
		UnitsPerFloor: tower.UnitsPerFloor(),
	}
}

// ToComplexSizeResponse converte um tamanho do domínio em DTO.
func ToComplexSizeResponse(size propertycoveragemodel.VerticalComplexSizeInterface) dto.ComplexSizeResponse {
	return dto.ComplexSizeResponse{
		ID:          size.ID(),
		ComplexID:   size.VerticalComplexID(),
		Size:        size.Size(),
		Description: size.Description(),
	}
}

// ToComplexZipCodeResponse converte um CEP do domínio em DTO.
func ToComplexZipCodeResponse(zip propertycoveragemodel.HorizontalComplexZipCodeInterface) dto.ComplexZipCodeResponse {
	return dto.ComplexZipCodeResponse{
		ID:        zip.ID(),
		ComplexID: zip.HorizontalComplexID(),
		ZipCode:   zip.ZipCode(),
	}
}

// ToListingComplexItems converts managed complexes to a minimal list payload.
func ToListingComplexItems(entities []propertycoveragemodel.ManagedComplexInterface) []dto.ListingComplexItemResponse {
	responses := make([]dto.ListingComplexItemResponse, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, dto.ListingComplexItemResponse{
			ID:   entity.ID(),
			Name: entity.Name(),
		})
	}
	return responses
}
