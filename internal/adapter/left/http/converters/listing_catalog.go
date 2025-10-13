package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// ToListingCatalogValueResponse converte um valor de catálogo do domínio para o DTO utilizado nas respostas HTTP.
func ToListingCatalogValueResponse(value listingmodel.CatalogValueInterface) dto.ListingCatalogValueResponse {
	resp := dto.ListingCatalogValueResponse{
		ID:       int(value.ID()),
		Category: value.Category(),
		Slug:     value.Slug(),
		Label:    value.Label(),
		IsActive: value.IsActive(),
	}

	if desc := value.Description(); desc != nil {
		resp.Description = desc
	}

	return resp
}

// ToListingCatalogValuesResponse converte uma coleção de valores de catálogo para o envelope de resposta HTTP.
func ToListingCatalogValuesResponse(values []listingmodel.CatalogValueInterface) dto.ListingCatalogValuesResponse {
	resp := dto.ListingCatalogValuesResponse{
		Values: make([]dto.ListingCatalogValueResponse, 0, len(values)),
	}

	for _, value := range values {
		resp.Values = append(resp.Values, ToListingCatalogValueResponse(value))
	}

	return resp
}
