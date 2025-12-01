package listinghandlers

import (
	listinghandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/listinghandler"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	listingservice "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
)

// ListingHandler implementa os handlers HTTP para operações de listing
type ListingHandler struct {
	listingService listingservice.ListingServiceInterface
	globalService  globalservice.GlobalServiceInterface
}

// NewListingHandlerAdapter cria uma nova instância de ListingHandler
func NewListingHandlerAdapter(
	listingService listingservice.ListingServiceInterface,
	globalService globalservice.GlobalServiceInterface,
) listinghandlerport.ListingHandlerPort {
	return &ListingHandler{
		listingService: listingService,
		globalService:  globalService,
	}
}
