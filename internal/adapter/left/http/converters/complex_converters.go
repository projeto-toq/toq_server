package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

// ToComplexResponse converte o domínio de empreendimento para DTO HTTP.
func ToComplexResponse(entity complexmodel.ComplexInterface) dto.ComplexResponse {
	resp := dto.ComplexResponse{
		ID:               entity.ID(),
		Name:             entity.Name(),
		ZipCode:          entity.ZipCode(),
		Street:           entity.Street(),
		Number:           entity.Number(),
		Neighborhood:     entity.Neighborhood(),
		City:             entity.City(),
		State:            entity.State(),
		PhoneNumber:      entity.PhoneNumber(),
		Sector:           uint8(entity.Sector()),
		MainRegistration: entity.MainRegistration(),
		PropertyType:     uint16(entity.GetPropertyType()),
	}

	if sizes := entity.ComplexSizes(); len(sizes) > 0 {
		resp.Sizes = make([]dto.ComplexSizeResponse, 0, len(sizes))
		for _, size := range sizes {
			resp.Sizes = append(resp.Sizes, ToComplexSizeResponse(size))
		}
	}

	if towers := entity.ComplexTowers(); len(towers) > 0 {
		resp.Towers = make([]dto.ComplexTowerResponse, 0, len(towers))
		for _, tower := range towers {
			resp.Towers = append(resp.Towers, ToComplexTowerResponse(tower))
		}
	}

	if zipCodes := entity.ComplexZipCodes(); len(zipCodes) > 0 {
		resp.ZipCodes = make([]dto.ComplexZipCodeResponse, 0, len(zipCodes))
		for _, zip := range zipCodes {
			resp.ZipCodes = append(resp.ZipCodes, ToComplexZipCodeResponse(zip))
		}
	}

	return resp
}

// ToComplexResponses converte vários empreendimentos para DTO.
func ToComplexResponses(entities []complexmodel.ComplexInterface) []dto.ComplexResponse {
	responses := make([]dto.ComplexResponse, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, ToComplexResponse(entity))
	}
	return responses
}

// ToComplexTowerResponse converte uma torre do domínio em DTO.
func ToComplexTowerResponse(tower complexmodel.ComplexTowerInterface) dto.ComplexTowerResponse {
	return dto.ComplexTowerResponse{
		ID:            tower.ID(),
		ComplexID:     tower.ComplexID(),
		Tower:         tower.Tower(),
		Floors:        tower.Floors(),
		TotalUnits:    tower.TotalUnits(),
		UnitsPerFloor: tower.UnitsPerFloor(),
	}
}

// ToComplexSizeResponse converte um tamanho do domínio em DTO.
func ToComplexSizeResponse(size complexmodel.ComplexSizeInterface) dto.ComplexSizeResponse {
	return dto.ComplexSizeResponse{
		ID:          size.ID(),
		ComplexID:   size.ComplexID(),
		Size:        size.Size(),
		Description: size.Description(),
	}
}

// ToComplexZipCodeResponse converte um CEP do domínio em DTO.
func ToComplexZipCodeResponse(zip complexmodel.ComplexZipCodeInterface) dto.ComplexZipCodeResponse {
	return dto.ComplexZipCodeResponse{
		ID:        zip.ID(),
		ComplexID: zip.ComplexID(),
		ZipCode:   zip.ZipCode(),
	}
}
