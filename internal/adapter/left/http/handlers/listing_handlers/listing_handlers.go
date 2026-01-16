package listinghandlers

import (
	listinghandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/listinghandler"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	listingservice "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
)

// ListingHandler implementa os handlers HTTP para operações de listing
type ListingHandler struct {
	listingService          listingservice.ListingServiceInterface
	globalService           globalservice.GlobalServiceInterface
	userService             userservices.UserServiceInterface
	propertyCoverageService propertycoverageservice.PropertyCoverageServiceInterface
	config                  ListingHandlerConfig
}

// ListingHandlerConfig carries HTTP-level configuration for listing handlers.
// Thresholds are resolved at bootstrap from environment values.
type ListingHandlerConfig struct {
	NewListingHoursThreshold   int
	PriceChangedHoursThreshold int
}

// NewListingHandlerAdapter cria uma nova instância de ListingHandler
func NewListingHandlerAdapter(
	listingService listingservice.ListingServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	userService userservices.UserServiceInterface,
	propertyCoverageService propertycoverageservice.PropertyCoverageServiceInterface,
	config ListingHandlerConfig,
) listinghandlerport.ListingHandlerPort {
	return &ListingHandler{
		listingService:          listingService,
		globalService:           globalService,
		userService:             userService,
		propertyCoverageService: propertyCoverageService,
		config:                  config,
	}
}
