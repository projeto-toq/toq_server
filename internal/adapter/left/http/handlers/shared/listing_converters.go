package shared

import (
	httpmodels "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/models"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
)

// ConvertDomainListingToHTTP converte um modelo de domínio de listing para o modelo HTTP
//
// Parâmetros:
//   - listing: Interface do modelo de domínio de listing
//
// Retorna:
//   - httpmodels.Listing: Modelo HTTP correspondente
func ConvertDomainListingToHTTP(listing listingmodel.ListingInterface) httpmodels.Listing {
	// Convert features - simplificado já que os métodos das features não correspondem
	httpFeatures := make([]httpmodels.BaseFeature, 0)
	// TODO: implementar conversão de features quando as interfaces estiverem alinhadas

	return httpmodels.Listing{
		ID:           listing.ID(),
		ZipCode:      listing.ZipCode(),
		Number:       listing.Number(),
		PropertyType: int32(listing.ListingType()),
		Features:     httpFeatures,
		Status:       listing.Status().String(),
		// Campos de timestamp serão implementados quando as interfaces forem atualizadas
		CreatedAt: "2023-01-01T00:00:00Z", // placeholder
		UpdatedAt: "2023-01-01T00:00:00Z", // placeholder
	}
}

// ConvertDomainOfferToHTTP converte um modelo de domínio de offer para o modelo HTTP
//
// Parâmetros:
//   - offer: Interface do modelo de domínio de offer
//
// Retorna:
//   - httpmodels.Offer: Modelo HTTP correspondente
func ConvertDomainOfferToHTTP(offer listingmodel.OfferInterface) httpmodels.Offer {
	return httpmodels.Offer{
		ID: offer.ID(),
		// Campos específicos serão implementados quando as interfaces forem expandidas
		ListingID: 0,                      // placeholder
		Amount:    0,                      // placeholder
		Status:    "todo",                 // placeholder
		CreatedAt: "2023-01-01T00:00:00Z", // placeholder
	}
}

// ConvertDomainVisitToHTTP converte um modelo de domínio de visit para o modelo HTTP
//
// Parâmetros:
//   - visit: Interface do modelo de domínio de visit
//
// Retorna:
//   - httpmodels.Visit: Modelo HTTP correspondente
func ConvertDomainVisitToHTTP(visit listingmodel.VisitInterface) httpmodels.Visit {
	return httpmodels.Visit{
		ID: visit.ID(),
		// Campos específicos serão implementados quando as interfaces forem expandidas
		ListingID: 0,                      // placeholder
		Date:      "2023-01-01T10:00:00Z", // placeholder
		Status:    "todo",                 // placeholder
		CreatedAt: "2023-01-01T00:00:00Z", // placeholder
	}
}
