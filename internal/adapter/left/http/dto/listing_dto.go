package dto

import (
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// Listing DTOs for HTTP handlers

// GetAllListingsRequest represents request for getting all listings
type GetAllListingsRequest struct {
	Page         int    `form:"page,default=1" binding:"min=1"`
	Limit        int    `form:"limit,default=10" binding:"min=1,max=100"`
	Status       string `form:"status,omitempty"`
	Code         string `form:"code,omitempty"`
	Title        string `form:"title,omitempty"`
	UserID       string `form:"userId,omitempty"`
	ZipCode      string `form:"zipCode,omitempty"`
	City         string `form:"city,omitempty"`
	Neighborhood string `form:"neighborhood,omitempty"`
	CreatedFrom  string `form:"createdFrom,omitempty"`
	CreatedTo    string `form:"createdTo,omitempty"`
	MinSellPrice string `form:"minSell,omitempty"`
	MaxSellPrice string `form:"maxSell,omitempty"`
	MinRentPrice string `form:"minRent,omitempty"`
	MaxRentPrice string `form:"maxRent,omitempty"`
	MinLandSize  string `form:"minLandSize,omitempty"`
	MaxLandSize  string `form:"maxLandSize,omitempty"`
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

// UpdateListingRequest represents request for updating a listing.
//
// Exemplo completo de payload (todos os campos preenchidos):
//
//	{
//	  "id": 98765,
//	  "owner": "myself",
//	  "features": [
//	    {"featureId": 101, "quantity": 2},
//	    {"featureId": 205, "quantity": 1}
//	  ],
//	  "landSize": 423.5,
//	  "corner": true,
//	  "nonBuildable": 12.75,
//	  "buildable": 410.75,
//	  "delivered": "furnished",
//	  "whoLives": "tenant",
//	  "description": "Apartamento amplo com vista panoramica",
//	  "transaction": "sale",
//	  "sellNet": 1200000,
//	  "rentNet": 8500,
//	  "condominium": 1200.5,
//	  "annualTax": 3400.75,
//	  "annualGroundRent": 1800,
//	  "exchange": true,
//	  "exchangePercentual": 50,
//	  "exchangePlaces": [
//	    {"neighborhood": "Vila Mariana", "city": "Sao Paulo", "state": "SP"},
//	    {"neighborhood": "Centro", "city": "Campinas", "state": "SP"}
//	  ],
//	  "installment": "short_term",
//	  "financing": true,
//	  "financingBlockers": ["pending_probate", "other"],
//	  "guarantees": [
//	    {"priority": 1, "guarantee": "security_deposit"},
//	    {"priority": 2, "guarantee": "surety_bond"}
//	  ],
//	  "visit": "client",
//	  "tenantName": "Joao da Silva",
//	  "tenantEmail": "joao.silva@example.com",
//	  "tenantPhone": "+55 11 91234-5678",
//	  "accompanying": "assistant"
//	}
type UpdateListingRequest struct {
	ID                 coreutils.Optional[int64]                               `json:"id"`
	Owner              coreutils.Optional[string]                              `json:"owner"`
	Features           coreutils.Optional[[]UpdateListingFeatureRequest]       `json:"features"`
	LandSize           coreutils.Optional[float64]                             `json:"landSize"`
	Corner             coreutils.Optional[bool]                                `json:"corner"`
	NonBuildable       coreutils.Optional[float64]                             `json:"nonBuildable"`
	Buildable          coreutils.Optional[float64]                             `json:"buildable"`
	Delivered          coreutils.Optional[string]                              `json:"delivered"`
	WhoLives           coreutils.Optional[string]                              `json:"whoLives"`
	Description        coreutils.Optional[string]                              `json:"description"`
	Transaction        coreutils.Optional[string]                              `json:"transaction"`
	SellNet            coreutils.Optional[float64]                             `json:"sellNet"`
	RentNet            coreutils.Optional[float64]                             `json:"rentNet"`
	Condominium        coreutils.Optional[float64]                             `json:"condominium"`
	AnnualTax          coreutils.Optional[float64]                             `json:"annualTax"`
	AnnualGroundRent   coreutils.Optional[float64]                             `json:"annualGroundRent"`
	Exchange           coreutils.Optional[bool]                                `json:"exchange"`
	ExchangePercentual coreutils.Optional[float64]                             `json:"exchangePercentual"`
	ExchangePlaces     coreutils.Optional[[]UpdateListingExchangePlaceRequest] `json:"exchangePlaces"`
	Installment        coreutils.Optional[string]                              `json:"installment"`
	Financing          coreutils.Optional[bool]                                `json:"financing"`
	FinancingBlockers  coreutils.Optional[[]string]                            `json:"financingBlockers"`
	Guarantees         coreutils.Optional[[]UpdateListingGuaranteeRequest]     `json:"guarantees"`
	Visit              coreutils.Optional[string]                              `json:"visit"`
	TenantName         coreutils.Optional[string]                              `json:"tenantName"`
	TenantEmail        coreutils.Optional[string]                              `json:"tenantEmail"`
	TenantPhone        coreutils.Optional[string]                              `json:"tenantPhone"`
	Accompanying       coreutils.Optional[string]                              `json:"accompanying"`
}

// UpdateListingResponse represents response for updating a listing
type UpdateListingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EndUpdateListingRequest representa o payload para finalizar a atualização de um listing.
type EndUpdateListingRequest struct {
	ListingID int64 `json:"listingId" binding:"required"`
}

// EndUpdateListingResponse representa a resposta ao finalizar a atualização de um listing.
type EndUpdateListingResponse struct {
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
	Priority  uint8  `json:"priority"`
	Guarantee string `json:"guarantee"`
}

// Listing catalog DTOs

// ListingCatalogValueResponse represents a catalog entry available for listings.
type ListingCatalogValueResponse struct {
	ID           int     `json:"id"`
	Category     string  `json:"category"`
	NumericValue int     `json:"numericValue"`
	Slug         string  `json:"slug"`
	Label        string  `json:"label"`
	Description  *string `json:"description,omitempty"`
	IsActive     bool    `json:"isActive"`
}

// ListingCatalogValuesResponse aggregates catalog values returned by the service.
type ListingCatalogValuesResponse struct {
	Values []ListingCatalogValueResponse `json:"values"`
}

// ListingCatalogRequest captures the payload used to list catalog values in the app.
type ListingCatalogRequest struct {
	Category string `json:"category" binding:"required"`
}

// AdminListingCatalogRequest defines the payload used to query the catalog from the admin panel.
type AdminListingCatalogRequest struct {
	Category        string `json:"category" binding:"required"`
	IncludeInactive bool   `json:"includeInactive"`
}

// ListingCatalogCreateRequest defines the payload for creating catalog values.
type ListingCatalogCreateRequest struct {
	Category    string                     `json:"category" binding:"required"`
	Slug        string                     `json:"slug" binding:"required"`
	Label       string                     `json:"label" binding:"required"`
	Description coreutils.Optional[string] `json:"description,omitempty"`
	IsActive    coreutils.Optional[bool]   `json:"isActive,omitempty"`
}

// ListingCatalogUpdateRequest defines the payload for partially updating catalog values.
type ListingCatalogUpdateRequest struct {
	ID          uint8                      `json:"id" binding:"required,min=1"`
	Category    string                     `json:"category" binding:"required"`
	Slug        coreutils.Optional[string] `json:"slug,omitempty"`
	Label       coreutils.Optional[string] `json:"label,omitempty"`
	Description coreutils.Optional[string] `json:"description,omitempty"`
	IsActive    coreutils.Optional[bool]   `json:"isActive,omitempty"`
}

// ListingCatalogDeleteRequest defines the payload for deactivating catalog values.
type ListingCatalogDeleteRequest struct {
	ID       uint8  `json:"id" binding:"required,min=1"`
	Category string `json:"category" binding:"required"`
}

// ListingCatalogRestoreRequest defines the payload for reactivating catalog values.
type ListingCatalogRestoreRequest struct {
	ID       uint8  `json:"id" binding:"required,min=1"`
	Category string `json:"category" binding:"required"`
}

// ListPhotographerSlotsRequest define filtros e paginação para consulta de slots.
type ListPhotographerSlotsRequest struct {
	From   string `form:"from" binding:"omitempty" example:"2025-10-20"`
	To     string `form:"to" binding:"omitempty" example:"2025-10-31"`
	Period string `form:"period" binding:"omitempty,oneof=MORNING AFTERNOON" example:"MORNING"`
	Page   int    `form:"page,default=1" binding:"min=1"`
	Size   int    `form:"size,default=20" binding:"min=1,max=100"`
	Sort   string `form:"sort,default=date_asc" binding:"omitempty,oneof=date_asc date_desc photographer_asc photographer_desc"`
}

// PhotographerSlotResponse representa um slot disponível na agenda dos fotógrafos.
type PhotographerSlotResponse struct {
	SlotID             uint64  `json:"slotId" example:"123"`
	PhotographerUserID uint64  `json:"photographerUserId" example:"45"`
	SlotDate           string  `json:"slotDate" example:"2025-10-25"`
	Period             string  `json:"period" example:"MORNING"`
	Status             string  `json:"status" example:"AVAILABLE"`
	ReservedUntil      *string `json:"reservedUntil,omitempty" example:"2025-10-24T12:00:00Z"`
}

// ListPhotographerSlotsResponse agrega slots e paginação.
type ListPhotographerSlotsResponse struct {
	Data       []PhotographerSlotResponse `json:"data"`
	Pagination PaginationResponse         `json:"pagination"`
}

// ReservePhotoSessionRequest representa o corpo para reservar um slot.
type ReservePhotoSessionRequest struct {
	ListingID int64  `json:"listingId" binding:"required" example:"1001"`
	SlotID    uint64 `json:"slotId" binding:"required" example:"2002"`
}

// ReservePhotoSessionResponse retorna dados da reserva temporária.
type ReservePhotoSessionResponse struct {
	SlotID           uint64 `json:"slotId" example:"2002"`
	ReservationToken string `json:"reservationToken" example:"c36b754f-6c37-4c15-8f25-9d77ddf9bb3e"`
	ExpiresAt        string `json:"expiresAt" example:"2025-10-24T14:45:00Z"`
}

// ConfirmPhotoSessionRequest representa o payload para confirmar a sessão de fotos.
type ConfirmPhotoSessionRequest struct {
	ListingID        int64  `json:"listingId" binding:"required" example:"1001"`
	SlotID           uint64 `json:"slotId" binding:"required" example:"2002"`
	ReservationToken string `json:"reservationToken" binding:"required" example:"c36b754f-6c37-4c15-8f25-9d77ddf9bb3e"`
}

// ConfirmPhotoSessionResponse retorna dados da sessão confirmada.
type ConfirmPhotoSessionResponse struct {
	PhotoSessionID uint64 `json:"photoSessionId" example:"3003"`
	SlotID         uint64 `json:"slotId" example:"2002"`
	ScheduledStart string `json:"scheduledStart" example:"2025-10-24T09:00:00Z"`
	ScheduledEnd   string `json:"scheduledEnd" example:"2025-10-24T13:00:00Z"`
}

// CancelPhotoSessionRequest representa o corpo para cancelamento de sessão de fotos.
type CancelPhotoSessionRequest struct {
	PhotoSessionID uint64 `json:"photoSessionId" binding:"required" example:"3003"`
}
