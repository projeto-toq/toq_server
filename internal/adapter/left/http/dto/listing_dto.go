package dto

import (
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// Listing DTOs for HTTP handlers

// GetAllListingsRequest represents request for getting all listings
type GetAllListingsRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=10" binding:"min=1,max=100"`
	Status   string `form:"status,omitempty"`
	UserID   string `form:"userId,omitempty"`
	ZipCode  string `form:"zipCode,omitempty"`
	MinPrice int    `form:"minPrice,omitempty"`
	MaxPrice int    `form:"maxPrice,omitempty"`
}

// GetAllListingsResponse represents response for getting all listings
type GetAllListingsResponse struct {
	Data       []ListingResponse  `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// StartListingRequest represents request for starting a new listing
type StartListingRequest struct {
	Number       string  `json:"number" binding:"required"`
	City         string  `json:"city" binding:"required"`
	Neighborhood *string `json:"neighborhood,omitempty"`
	Complement   *string `json:"complement,omitempty"`
	State        string  `json:"state" binding:"required"`
	Street       string  `json:"street" binding:"required"`
	ZipCode      string  `json:"zipCode" binding:"required"`
	PropertyType int     `json:"propertyType" binding:"required"`
}

// StartListingResponse represents response for starting a new listing
type StartListingResponse struct {
	ID int64 `json:"id"`
}

// UpdateListingRequest represents request for updating a listing
type UpdateListingRequest struct {
	ID                 coreutils.Optional[int64]                               `json:"id"`
	Owner              coreutils.Optional[int]                                 `json:"owner"`
	Features           coreutils.Optional[[]UpdateListingFeatureRequest]       `json:"features"`
	LandSize           coreutils.Optional[float64]                             `json:"landSize"`
	Corner             coreutils.Optional[bool]                                `json:"corner"`
	NonBuildable       coreutils.Optional[float64]                             `json:"nonBuildable"`
	Buildable          coreutils.Optional[float64]                             `json:"buildable"`
	Delivered          coreutils.Optional[int]                                 `json:"delivered"`
	WhoLives           coreutils.Optional[int]                                 `json:"whoLives"`
	Description        coreutils.Optional[string]                              `json:"description"`
	Transaction        coreutils.Optional[int]                                 `json:"transaction"`
	SellNet            coreutils.Optional[float64]                             `json:"sellNet"`
	RentNet            coreutils.Optional[float64]                             `json:"rentNet"`
	Condominium        coreutils.Optional[float64]                             `json:"condominium"`
	AnnualTax          coreutils.Optional[float64]                             `json:"annualTax"`
	AnnualGroundRent   coreutils.Optional[float64]                             `json:"annualGroundRent"`
	Exchange           coreutils.Optional[bool]                                `json:"exchange"`
	ExchangePercentual coreutils.Optional[float64]                             `json:"exchangePercentual"`
	ExchangePlaces     coreutils.Optional[[]UpdateListingExchangePlaceRequest] `json:"exchangePlaces"`
	Installment        coreutils.Optional[int]                                 `json:"installment"`
	Financing          coreutils.Optional[bool]                                `json:"financing"`
	FinancingBlockers  coreutils.Optional[[]int]                               `json:"financingBlockers"`
	Guarantees         coreutils.Optional[[]UpdateListingGuaranteeRequest]     `json:"guarantees"`
	Visit              coreutils.Optional[int]                                 `json:"visit"`
	TenantName         coreutils.Optional[string]                              `json:"tenantName"`
	TenantEmail        coreutils.Optional[string]                              `json:"tenantEmail"`
	TenantPhone        coreutils.Optional[string]                              `json:"tenantPhone"`
	Accompanying       coreutils.Optional[int]                                 `json:"accompanying"`
}

// UpdateListingResponse represents response for updating a listing
type UpdateListingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DeleteListingResponse represents response for deleting a listing
type DeleteListingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ChangeListingStatusRequest represents request for changing listing status
type ChangeListingStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// ChangeListingStatusResponse represents response for changing listing status
type ChangeListingStatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetOptionsResponse represents response for getting options
type GetOptionsResponse struct {
	PropertyTypes []PropertyTypeOption `json:"propertyTypes"`
}

// GetOptionsRequest representa o payload para obter opções de listing
type GetOptionsRequest struct {
	ZipCode string `json:"zipCode" binding:"required"`
	Number  string `json:"number" binding:"required"`
}

// PropertyTypeOption represents a property type option
type PropertyTypeOption struct {
	PropertyType int    `json:"propertyType"`
	Name         string `json:"name"`
}

// GetBaseFeaturesResponse represents response for getting base features
type GetBaseFeaturesResponse struct {
	Features []BaseFeature `json:"features"`
}

// BaseFeature represents a base feature
type BaseFeature struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// ListingResponse represents a listing in responses
type ListingResponse struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	PropertyType int     `json:"propertyType"`
	ZipCode      string  `json:"zipCode"`
	Number       string  `json:"number"`
	UserID       int64   `json:"userId"`
	ComplexID    string  `json:"complexId,omitempty"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

// AddListingPhotosRequest represents request for adding photos to a listing
type AddListingPhotosRequest struct {
	Photos []PhotoRequest `json:"photos" binding:"required,min=1"`
}

// PhotoRequest represents a photo in requests
type PhotoRequest struct {
	URL         string `json:"url" binding:"required"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order,omitempty"`
}

// AddListingPhotosResponse represents response for adding photos
type AddListingPhotosResponse struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	PhotoIDs []string `json:"photoIds,omitempty"`
}

// UpdateListingPhotoRequest represents request for updating a photo
type UpdateListingPhotoRequest struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order,omitempty"`
}

// UpdateListingPhotoResponse represents response for updating a photo
type UpdateListingPhotoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RemoveListingPhotoResponse represents response for removing a photo
type RemoveListingPhotoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// NotImplementedResponse represents response for not implemented endpoints
type NotImplementedResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// UpdateListingFeatureRequest representa uma feature enviada no update de listing.
type UpdateListingFeatureRequest struct {
	FeatureID int64 `json:"featureId"`
	Quantity  uint8 `json:"quantity"`
}

// UpdateListingExchangePlaceRequest representa uma localidade de troca no payload de update de listing.
type UpdateListingExchangePlaceRequest struct {
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

// UpdateListingGuaranteeRequest representa uma garantia enviada no update de listing.
type UpdateListingGuaranteeRequest struct {
	Priority  uint8 `json:"priority"`
	Guarantee int   `json:"guarantee"`
}
