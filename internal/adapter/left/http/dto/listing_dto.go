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
	ZipCode      string  `json:"zipCode" binding:"required" example:"06543001" description:"Zip code without separators (8 digits)."`
	PropertyType int     `json:"propertyType" binding:"required"`
}

// StartListingResponse represents response for starting a new listing
type StartListingResponse struct {
	ID                int64  `json:"id"`
	ListingIdentityID int64  `json:"listingIdentityId"`
	ListingUUID       string `json:"listingUuid"`
	ActiveVersionID   int64  `json:"activeVersionId"`
	Version           uint8  `json:"version"`
	Status            string `json:"status"`
}

// CreateDraftVersionRequest represents request for creating a draft version from active listing
type CreateDraftVersionRequest struct {
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required"`
}

// CreateDraftVersionResponse represents response for draft version creation
type CreateDraftVersionResponse struct {
	VersionID int64  `json:"versionId"`
	Version   uint8  `json:"version"`
	Status    string `json:"status"`
}

// GetListingDetailRequest representa a requisição para obter os detalhes completos de um listing.
type GetListingDetailRequest struct {
	ListingID int64 `json:"listingId" binding:"required"`
}

// CatalogItemResponse padroniza valores de catálogo retornados pelo endpoint de detalhes de listing.
type CatalogItemResponse struct {
	NumericValue uint8  `json:"numericValue"`
	Slug         string `json:"slug,omitempty"`
	Label        string `json:"label,omitempty"`
}

// ListingPropertyTypeResponse expõe metadados do tipo de imóvel associado ao listing.
type ListingPropertyTypeResponse struct {
	Code        int64  `json:"code,omitempty"`
	Label       string `json:"label,omitempty"`
	PropertyBit uint16 `json:"propertyBit"`
}

// ListingFeatureResponse representa uma feature do listing com metadados enriquecidos.
type ListingFeatureResponse struct {
	Feature     string `json:"feature"`
	Description string `json:"description,omitempty"`
	Quantity    uint8  `json:"quantity"`
}

// ListingExchangePlaceResponse expõe um local aceito para permuta.
type ListingExchangePlaceResponse struct {
	ID               int64  `json:"id"`
	ListingID        int64  `json:"listingId"`
	ListingVersionID int64  `json:"listingVersionId"`
	Neighborhood     string `json:"neighborhood"`
	City             string `json:"city"`
	State            string `json:"state"`
}

// ListingFinancingBlockerResponse representa um impeditivo de financiamento associado ao listing.
type ListingFinancingBlockerResponse struct {
	ID               int64               `json:"id"`
	ListingID        int64               `json:"listingId"`
	ListingVersionID int64               `json:"listingVersionId"`
	Blocker          CatalogItemResponse `json:"blocker"`
}

// ListingGuaranteeResponse representa uma garantia aceita no listing.
type ListingGuaranteeResponse struct {
	ID               int64               `json:"id"`
	ListingID        int64               `json:"listingId"`
	ListingVersionID int64               `json:"listingVersionId"`
	Priority         uint8               `json:"priority"`
	Guarantee        CatalogItemResponse `json:"guarantee"`
}

// ListingDetailResponse agrega todos os campos do listing.
type ListingDetailResponse struct {
	ID                 int64                             `json:"id"`
	ListingIdentityID  int64                             `json:"listingIdentityId"`
	ListingUUID        string                            `json:"listingUuid"`
	ActiveVersionID    int64                             `json:"activeVersionId"`
	DraftVersionID     *int64                            `json:"draftVersionId,omitempty"`
	UserID             int64                             `json:"userId"`
	Code               uint32                            `json:"code"`
	Version            uint8                             `json:"version"`
	Status             string                            `json:"status"`
	ZipCode            string                            `json:"zipCode"`
	Street             string                            `json:"street"`
	Number             string                            `json:"number"`
	Complement         string                            `json:"complement"`
	Neighborhood       string                            `json:"neighborhood"`
	City               string                            `json:"city"`
	State              string                            `json:"state"`
	Title              string                            `json:"title"`
	PropertyType       *ListingPropertyTypeResponse      `json:"propertyType,omitempty"`
	Owner              *CatalogItemResponse              `json:"owner,omitempty"`
	Features           []ListingFeatureResponse          `json:"features,omitempty"`
	LandSize           float64                           `json:"landSize"`
	Corner             bool                              `json:"corner"`
	NonBuildable       float64                           `json:"nonBuildable"`
	Buildable          float64                           `json:"buildable"`
	Delivered          *CatalogItemResponse              `json:"delivered,omitempty"`
	WhoLives           *CatalogItemResponse              `json:"whoLives,omitempty"`
	Description        string                            `json:"description"`
	Transaction        *CatalogItemResponse              `json:"transaction,omitempty"`
	SellNet            float64                           `json:"sellNet"`
	RentNet            float64                           `json:"rentNet"`
	Condominium        float64                           `json:"condominium"`
	AnnualTax          float64                           `json:"annualTax"`
	MonthlyTax         float64                           `json:"monthlyTax"`
	AnnualGroundRent   float64                           `json:"annualGroundRent"`
	MonthlyGroundRent  float64                           `json:"monthlyGroundRent"`
	Exchange           bool                              `json:"exchange"`
	ExchangePercentual float64                           `json:"exchangePercentual"`
	ExchangePlaces     []ListingExchangePlaceResponse    `json:"exchangePlaces,omitempty"`
	Installment        *CatalogItemResponse              `json:"installment,omitempty"`
	Financing          bool                              `json:"financing"`
	FinancingBlockers  []ListingFinancingBlockerResponse `json:"financingBlockers,omitempty"`
	Guarantees         []ListingGuaranteeResponse        `json:"guarantees,omitempty"`
	Visit              *CatalogItemResponse              `json:"visit,omitempty"`
	TenantName         string                            `json:"tenantName"`
	TenantEmail        string                            `json:"tenantEmail"`
	TenantPhone        string                            `json:"tenantPhone"`
	Accompanying       *CatalogItemResponse              `json:"accompanying,omitempty"`
	PhotoSessionID     *uint64                           `json:"photoSessionId,omitempty"`
	Deleted            bool                              `json:"deleted"`
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
//	  "tenantPhone": "+5511912345678",
//	  "accompanying": "assistant"
//	}
type UpdateListingRequest struct {
	ListingIdentityID  coreutils.Optional[int64]                               `json:"listingIdentityId" binding:"required" example:"1024"`
	ListingVersionID   coreutils.Optional[int64]                               `json:"listingVersionId" binding:"required" example:"5001"`
	Owner              coreutils.Optional[string]                              `json:"owner"`
	Features           coreutils.Optional[[]UpdateListingFeatureRequest]       `json:"features"`
	LandSize           coreutils.Optional[float64]                             `json:"landSize"`
	Corner             coreutils.Optional[bool]                                `json:"corner"`
	NonBuildable       coreutils.Optional[float64]                             `json:"nonBuildable"`
	Buildable          coreutils.Optional[float64]                             `json:"buildable"`
	Delivered          coreutils.Optional[string]                              `json:"delivered"`
	WhoLives           coreutils.Optional[string]                              `json:"whoLives"`
	Title              coreutils.Optional[string]                              `json:"title"`
	Description        coreutils.Optional[string]                              `json:"description"`
	Transaction        coreutils.Optional[string]                              `json:"transaction"`
	SellNet            coreutils.Optional[float64]                             `json:"sellNet"`
	RentNet            coreutils.Optional[float64]                             `json:"rentNet"`
	Condominium        coreutils.Optional[float64]                             `json:"condominium"`
	AnnualTax          coreutils.Optional[float64]                             `json:"annualTax" description:"Annual IPTU (property tax). Mutually exclusive with monthlyTax."`
	MonthlyTax         coreutils.Optional[float64]                             `json:"monthlyTax" description:"Monthly IPTU (property tax). Mutually exclusive with annualTax."`
	AnnualGroundRent   coreutils.Optional[float64]                             `json:"annualGroundRent" description:"Annual Laudêmio (ground rent). Mutually exclusive with monthlyGroundRent."`
	MonthlyGroundRent  coreutils.Optional[float64]                             `json:"monthlyGroundRent" description:"Monthly Laudêmio (ground rent). Mutually exclusive with annualGroundRent."`
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
	TenantPhone        coreutils.Optional[string]                              `json:"tenantPhone" description:"Tenant phone number in E.164 format (e.g., +5511912345678)."`
	Accompanying       coreutils.Optional[string]                              `json:"accompanying"`
}

// UpdateListingResponse represents response for updating a listing
type UpdateListingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PromoteListingVersionRequest encapsula a versão a ser promovida para ativa.
type PromoteListingVersionRequest struct {
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`
	VersionID         int64 `json:"versionId" binding:"required,min=1" example:"5001"`
}

// PromoteListingVersionResponse confirma a promoção de versão.
type PromoteListingVersionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DiscardDraftVersionRequest identifica o draft que deve ser descartado.
type DiscardDraftVersionRequest struct {
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`
	VersionID         int64 `json:"versionId" binding:"required,min=1" example:"5001"`
}

// DiscardDraftVersionResponse confirma o descarte de um draft de anúncio.
type DiscardDraftVersionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ListListingVersionsRequest captures filters for querying versions of a listing identity
//
// This DTO is used to retrieve all version history for a specific listing.
// The listing identity ID is mandatory, while the includeDeleted flag is optional.
//
// Business Rules:
//   - listingIdentityId must be a valid, existing listing identity ID
//   - includeDeleted defaults to false (only active versions returned)
//   - Returns all versions ordered by version number descending
type ListListingVersionsRequest struct {
	// ListingIdentityID is the unique identifier of the listing identity
	// Required field to identify which listing's versions should be retrieved
	// Must be greater than 0
	// Example: 1024
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`

	// IncludeDeleted determines whether soft-deleted versions should be included in results
	// When true, returns both active and soft-deleted versions
	// When false (default), returns only active versions (deleted = 0)
	// Example: false
	IncludeDeleted bool `json:"includeDeleted" binding:"omitempty" example:"false"`
}

// ListingVersionSummaryResponse exposes metadata for a specific listing version
//
// Represents a single version in the listing's history, including its status and active flag.
type ListingVersionSummaryResponse struct {
	// ID is the unique identifier for this specific version
	ID int64 `json:"id" example:"5001"`

	// ListingIdentityID is the parent listing identity this version belongs to
	ListingIdentityID int64 `json:"listingIdentityId" example:"1024"`

	// ListingUUID is the immutable UUID for the listing identity
	ListingUUID string `json:"listingUuid" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Version is the sequential version number (1, 2, 3, ...)
	Version uint8 `json:"version" example:"2"`

	// Status indicates the current lifecycle state of this version
	// Possible values: "DRAFT", "ACTIVE", "INACTIVE", "DELETED"
	Status string `json:"status" example:"ACTIVE"`

	// Title is the listing title for this version (may be empty for drafts)
	Title string `json:"title,omitempty" example:"Apartamento 3 dormitórios na Vila Mariana"`

	// IsActive indicates if this version is currently the active one
	// Only one version per listing identity should have IsActive = true
	IsActive bool `json:"isActive" example:"true"`
}

// ListListingVersionsResponse aggregates version metadata for a listing identity
//
// Returns complete version history with the active version flagged.
// Versions are ordered by version number descending (newest first).
type ListListingVersionsResponse struct {
	// ListingIdentityID is the identity whose versions are being returned
	ListingIdentityID int64 `json:"listingIdentityId" example:"1024"`

	// ListingUUID is the immutable UUID for this listing identity
	ListingUUID string `json:"listingUuid" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Versions contains the list of all versions for this listing
	// Array is ordered by version number descending
	Versions []ListingVersionSummaryResponse `json:"versions"`
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
	ID                int64   `json:"id"`
	ListingIdentityID int64   `json:"listingIdentityId"`
	ListingUUID       string  `json:"listingUuid"`
	ActiveVersionID   int64   `json:"activeVersionId"`
	DraftVersionID    *int64  `json:"draftVersionId,omitempty"`
	Version           uint8   `json:"version"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Price             float64 `json:"price"`
	Status            string  `json:"status"`
	PropertyType      int     `json:"propertyType"`
	ZipCode           string  `json:"zipCode"`
	Number            string  `json:"number"`
	UserID            int64   `json:"userId"`
	ComplexID         string  `json:"complexId,omitempty"`
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
	Category        string `json:"category" form:"category" binding:"required"`
	IncludeInactive bool   `json:"includeInactive" form:"includeInactive"`
}

// AdminGetListingCatalogDetailRequest represents the payload to fetch a catalog value detail.
type AdminGetListingCatalogDetailRequest struct {
	Category string `json:"category" form:"category" binding:"required"`
	ID       int    `json:"id" form:"id" binding:"required,min=1,max=255"`
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
	From              string `form:"from" binding:"omitempty" example:"2025-10-20"`
	To                string `form:"to" binding:"omitempty" example:"2025-10-31"`
	Period            string `form:"period" binding:"omitempty,oneof=MORNING AFTERNOON" example:"MORNING"`
	Page              int    `form:"page,default=1" binding:"min=1"`
	Size              int    `form:"size,default=20" binding:"min=1,max=100"`
	Sort              string `form:"sort,default=start_asc" binding:"omitempty,oneof=start_asc start_desc photographer_asc photographer_desc date_asc date_desc"`
	ListingIdentityID int64  `form:"listingIdentityId" binding:"required,min=1" example:"1024"`
	Timezone          string `form:"timezone" binding:"required" example:"America/Sao_Paulo"`
}

// PhotographerSlotResponse representa um slot disponível na agenda dos fotógrafos.
type PhotographerSlotResponse struct {
	SlotID             uint64  `json:"slotId" example:"123"`
	PhotographerUserID uint64  `json:"photographerUserId" example:"45"`
	SlotStart          string  `json:"slotStart" example:"2025-10-25T09:00:00Z"`
	SlotEnd            string  `json:"slotEnd" example:"2025-10-25T10:00:00Z"`
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
	ListingIdentityID int64  `json:"listingIdentityId" binding:"required" example:"1024"`
	SlotID            uint64 `json:"slotId" binding:"required" example:"2002"`
}

// ReservePhotoSessionResponse retorna dados da reserva temporária.
type ReservePhotoSessionResponse struct {
	SlotID         uint64 `json:"slotId" example:"2002"`
	SlotStart      string `json:"slotStart" example:"2025-10-24T09:00:00Z"`
	SlotEnd        string `json:"slotEnd" example:"2025-10-24T10:00:00Z"`
	PhotoSessionID uint64 `json:"photoSessionId" example:"3003"`
}

// ConfirmPhotoSessionRequest representa o payload para confirmar a sessão de fotos.
type ConfirmPhotoSessionRequest struct {
	ListingID      int64  `json:"listingId" binding:"required" example:"1001"`
	PhotoSessionID uint64 `json:"photoSessionId" binding:"required" example:"3003"`
}

// ConfirmPhotoSessionResponse retorna dados da sessão confirmada.
type ConfirmPhotoSessionResponse struct {
	PhotoSessionID uint64 `json:"photoSessionId" example:"3003"`
	ScheduledStart string `json:"scheduledStart" example:"2025-10-24T09:00:00Z"`
	ScheduledEnd   string `json:"scheduledEnd" example:"2025-10-24T10:00:00Z"`
	Status         string `json:"status" example:"ACTIVE"`
}

// CancelPhotoSessionRequest representa o corpo para cancelamento de sessão de fotos.
type CancelPhotoSessionRequest struct {
	PhotoSessionID uint64 `json:"photoSessionId" binding:"required" example:"3003"`
}
