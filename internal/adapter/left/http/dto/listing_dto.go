package dto

import (
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// Listing DTOs for HTTP handlers

// ListListingsRequest captures filters, pagination, and sorting for listing search
//
// This DTO is used to retrieve paginated, filtered, and sorted lists of active listing versions.
// By default, only versions linked as active_version_id in listing_identities are returned.
// Set includeAllVersions=true to retrieve all versions regardless of active status.
//
// Sorting:
//   - sortBy: Field to order results by (id, status)
//   - sortOrder: Direction (asc or desc)
//   - Default: id DESC (newest listings first)
//
// Filtering:
//   - status: Listing lifecycle status (e.g., PUBLISHED, DRAFT, UNDER_OFFER)
//   - code: Exact listing code match
//   - title: Wildcard search on title/description (supports '*')
//   - userId: Owner filter (auto-enforced for owner role)
//   - Location: zipCode, city, neighborhood, street, number, complement, complex, state (wildcards supported)
//   - Price ranges: minSell/maxSell, minRent/maxRent
//   - Size range: minLandSize/maxLandSize
//
// Pagination:
//   - page: 1-indexed page number
//   - limit: Items per page (max 100)
type ListListingsRequest struct {
	// Page number for pagination (1-indexed)
	// Minimum: 1, Default: 1
	// Example: 1
	Page int `form:"page,default=1" binding:"min=1" example:"1"`

	// Limit is the number of items per page
	// Minimum: 1, Maximum: 100, Default: 20
	// Example: 20
	Limit int `form:"limit,default=20" binding:"min=1,max=100" example:"20"`

	// SortBy specifies the field to order results by
	// Allowed values: id, status, zipCode, city, neighborhood, street, number, state, complex
	// Default: id (creation order proxy - higher ID = newer listing)
	// Example: "id"
	SortBy string `form:"sortBy,default=id" binding:"omitempty,oneof=id status zipCode city neighborhood street number state complex" example:"id"`

	// SortOrder specifies the sort direction
	// Allowed values: asc (ascending), desc (descending)
	// Default: desc (newest first)
	// Example: "desc"
	SortOrder string `form:"sortOrder,default=desc" binding:"omitempty,oneof=asc desc" example:"desc"`

	// Status filters by listing lifecycle status
	// Accepts enum name (e.g., "PUBLISHED") or numeric value
	// Optional - omit to retrieve all statuses
	// Example: "PUBLISHED"
	Status string `form:"status,omitempty" example:"PUBLISHED"`

	// Code filters by exact listing code
	// Optional - omit to retrieve all codes
	// Example: 1024
	Code string `form:"code,omitempty" example:"1024"`

	// Title performs wildcard search on listing title and description
	// Supports '*' as wildcard character
	// Case-insensitive partial match
	// Optional
	// Example: "*garden*"
	Title string `form:"title,omitempty" example:"*garden*"`

	// UserID filters listings by owner user ID
	// For owner role, this is auto-enforced to requester's ID (security)
	// Optional for admin/realtor roles
	// Example: 55
	UserID string `form:"userId,omitempty" example:"55"`

	// ZipCode filters by Brazilian postal code (CEP)
	// Digits only, supports '*' wildcard
	// Example: "06543*"
	ZipCode string `form:"zipCode,omitempty" example:"06543*"`

	// City filters by city name
	// Supports '*' wildcard for partial match
	// Case-insensitive
	// Example: "*Paulista*"
	City string `form:"city,omitempty" example:"*Paulista*"`

	// Neighborhood filters by neighborhood name
	// Supports '*' wildcard for partial match
	// Case-insensitive
	// Example: "*Centro*"
	Neighborhood string `form:"neighborhood,omitempty" example:"*Centro*"`

	// Street filters by street name (supports wildcard)
	// Example: "*Paulista*"
	Street string `form:"street,omitempty" example:"*Paulista*"`

	// Number filters by address number (allows wildcard for ranges and "S/N" cases)
	// Example: "12*"
	Number string `form:"number,omitempty" example:"12*"`

	// Complement filters by address complement (apartment, block, etc.)
	// Example: "*Bloco B*"
	Complement string `form:"complement,omitempty" example:"*Bloco B*"`

	// Complex filters by condominium/complex name (supports wildcard)
	// Example: "*Residencial Atlântico*"
	Complex string `form:"complex,omitempty" example:"*Residencial Atlântico*"`

	// State filters by Brazilian state (UF). Accepts wildcard but recommended to use exact 2-letter code
	// Example: "SP"
	State string `form:"state,omitempty" example:"SP"`

	// MinSellPrice filters listings with sell price >= this value
	// Optional - used with maxSell to define price range
	// Example: 100000
	MinSellPrice string `form:"minSell,omitempty" example:"100000"`

	// MaxSellPrice filters listings with sell price <= this value
	// Optional - used with minSell to define price range
	// Example: 900000
	MaxSellPrice string `form:"maxSell,omitempty" example:"900000"`

	// MinRentPrice filters listings with rent price >= this value
	// Optional - used with maxRent to define price range
	// Example: 1500
	MinRentPrice string `form:"minRent,omitempty" example:"1500"`

	// MaxRentPrice filters listings with rent price <= this value
	// Optional - used with minRent to define price range
	// Example: 8000
	MaxRentPrice string `form:"maxRent,omitempty" example:"8000"`

	// MinLandSize filters listings with land size >= this value (square meters)
	// Optional - used with maxLandSize to define size range
	// Example: 120.5
	MinLandSize string `form:"minLandSize,omitempty" example:"120.5"`

	// MaxLandSize filters listings with land size <= this value (square meters)
	// Optional - used with minLandSize to define size range
	// Example: 500.75
	MaxLandSize string `form:"maxLandSize,omitempty" example:"500.75"`

	// MinSuites filters listings whose suite count (feature "Suites") is >= this value
	// Optional - used with maxSuites to define suite range
	// Example: 2
	MinSuites string `form:"minSuites,omitempty" example:"2"`

	// MaxSuites filters listings whose suite count (feature "Suites") is <= this value
	// Optional - used with minSuites to define suite range
	// Example: 4
	MaxSuites string `form:"maxSuites,omitempty" example:"4"`

	// IncludeAllVersions determines version filtering behavior
	// false (default): Only active versions (linked via active_version_id)
	// true: All versions (active + draft)
	// Example: false
	IncludeAllVersions bool `form:"includeAllVersions,default=false" example:"false"`
}

// ListListingsResponse aggregates paginated listing data with metadata
//
// Contains the filtered and sorted listing collection plus pagination info.
type ListListingsResponse struct {
	// Data contains the listing collection for the current page
	Data []ListingResponse `json:"data"`

	// Pagination provides metadata for navigating the result set
	Pagination PaginationResponse `json:"pagination"`
}

// GetAllListingsRequest represents request for getting all listings (DEPRECATED - use ListListingsRequest)
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

// GetAllListingsResponse represents response for getting all listings (DEPRECATED - use ListListingsResponse)
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

	// New fields for duplicity check
	Complex    *string `json:"complex,omitempty"`
	UnitTower  *string `json:"unitTower,omitempty"`
	UnitFloor  *int16  `json:"unitFloor,omitempty"`
	UnitNumber *string `json:"unitNumber,omitempty"`
	LandBlock  *string `json:"landBlock,omitempty"`
	LandLot    *string `json:"landLot,omitempty"`
}

// StartListingResponse represents response for starting a new listing
//
// @Description Response after creating a new listing. Use listingIdentityId and listingVersionId for subsequent PUT /listing updates.
type StartListingResponse struct {
	ListingVersionID  int64  `json:"listingVersionId" example:"5001"`
	ListingIdentityID int64  `json:"listingIdentityId" example:"1024"`
	ListingUUID       string `json:"listingUuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Version           uint8  `json:"version" example:"1"`
	Status            string `json:"status" example:"DRAFT"`
}

// ChangeListingStatusRequest represents owner-driven transitions between READY and PUBLISHED states.
//
// This DTO is used by POST /listings/status and requires authentication. Only the listing owner can
// publish or suspend a listing, and the service validates both ownership and status transitions.
type ChangeListingStatusRequest struct {
	// ListingIdentityID identifies the listing identity (listing_identities.id) that will change status.
	// Must be greater than zero and belong to the authenticated owner.
	// Example: 1024
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`

	// Action determines the transition to apply.
	// Allowed values: PUBLISH (READY → PUBLISHED) or SUSPEND (PUBLISHED/UNDER_OFFER/UNDER_NEGOTIATION → READY).
	// Example: "PUBLISH"
	Action string `json:"action" binding:"required,oneof=PUBLISH SUSPEND" example:"PUBLISH"`
}

// ChangeListingStatusResponse echoes the updated metadata after applying the transition.
type ChangeListingStatusResponse struct {
	// ListingIdentityID is the identity that had its active version updated.
	ListingIdentityID int64 `json:"listingIdentityId" example:"1024"`

	// ActiveVersionID is the version affected by the transition.
	ActiveVersionID int64 `json:"activeVersionId" example:"5003"`

	// PreviousStatus is the status before the change.
	PreviousStatus string `json:"previousStatus" example:"READY"`

	// NewStatus is the resulting status after the action.
	NewStatus string `json:"newStatus" example:"PUBLISHED"`
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

// GetListingDetailRequest represents the request payload to retrieve full listing details
//
// This DTO is used to fetch comprehensive information about a listing, including:
//   - Active version data (the currently promoted version)
//   - Draft version metadata (if exists and requester is the owner)
//   - All property-specific fields and related entities (features, guarantees, exchange places, etc.)
//   - Photo session booking status
//
// The endpoint returns the ACTIVE version by default (referenced by listing_identities.active_version_id).
// Draft versions are exposed via draftVersionId field if they exist.
//
// Authorization:
//   - Only the listing owner (listing_identities.user_id == authenticated user_id) can access details
//   - Returns 403 Forbidden if requester is not the owner
//
// Example:
//
//	Request Body: {"listingIdentityId": 1024}
//	Response includes: active version fields + draftVersionId: 5002 (if draft exists)
type GetListingDetailRequest struct {
	// ListingIdentityID is the unique identifier of the listing identity (listing_identities.id)
	// This references the parent listing entity, not a specific version
	// Required field to identify which listing's details should be retrieved
	// Must be greater than 0
	// Example: 1024
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`
}

// FavoriteListingRequest represents the payload to add or remove a favorite.
type FavoriteListingRequest struct {
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required,min=1"`
}

// ListFavoriteListingsRequest captures pagination for favorites listing endpoint.
type ListFavoriteListingsRequest struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
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

// ListingOwnerInfoResponse expõe dados do proprietário e seus tempos de resposta.
type ListingOwnerInfoResponse struct {
	ID                     int64  `json:"id"`
	FullName               string `json:"fullName"`
	PhotoURL               string `json:"photoUrl,omitempty"`
	MemberSinceMonths      int    `json:"memberSinceMonths"`
	VisitAverageSeconds    *int64 `json:"visitAverageSeconds,omitempty"`
	ProposalAverageSeconds *int64 `json:"proposalAverageSeconds,omitempty"`
}

// ListingPerformanceMetricsResponse agrega métricas de engajamento do imóvel.
type ListingPerformanceMetricsResponse struct {
	Shares    int64 `json:"shares"`
	Views     int64 `json:"views"`
	Favorites int64 `json:"favorites"`
}

// ListingDetailResponse agrega todos os campos do listing.
type ListingDetailResponse struct {
	ID                         int64                             `json:"id"`
	ListingIdentityID          int64                             `json:"listingIdentityId"`
	ListingUUID                string                            `json:"listingUuid"`
	ActiveVersionID            int64                             `json:"activeVersionId"`
	DraftVersionID             *int64                            `json:"draftVersionId,omitempty"`
	UserID                     int64                             `json:"userId"`
	Code                       uint32                            `json:"code"`
	Version                    uint8                             `json:"version"`
	Status                     string                            `json:"status"`
	ZipCode                    string                            `json:"zipCode"`
	Street                     string                            `json:"street"`
	Number                     string                            `json:"number"`
	Complement                 string                            `json:"complement"`
	Neighborhood               string                            `json:"neighborhood"`
	City                       string                            `json:"city"`
	State                      string                            `json:"state"`
	Complex                    string                            `json:"complex,omitempty"`
	Title                      string                            `json:"title"`
	PropertyType               *ListingPropertyTypeResponse      `json:"propertyType,omitempty"`
	Owner                      *CatalogItemResponse              `json:"owner,omitempty"`
	OwnerInfo                  *ListingOwnerInfoResponse         `json:"ownerInfo,omitempty"`
	Features                   []ListingFeatureResponse          `json:"features,omitempty"`
	LandSize                   float64                           `json:"landSize"`
	Corner                     bool                              `json:"corner"`
	NonBuildable               float64                           `json:"nonBuildable"`
	Buildable                  float64                           `json:"buildable"`
	Delivered                  *CatalogItemResponse              `json:"delivered,omitempty"`
	WhoLives                   *CatalogItemResponse              `json:"whoLives,omitempty"`
	Description                string                            `json:"description"`
	Transaction                *CatalogItemResponse              `json:"transaction,omitempty"`
	SellNet                    float64                           `json:"sellNet"`
	RentNet                    float64                           `json:"rentNet"`
	Condominium                float64                           `json:"condominium"`
	AnnualTax                  float64                           `json:"annualTax"`
	MonthlyTax                 float64                           `json:"monthlyTax"`
	AnnualGroundRent           float64                           `json:"annualGroundRent"`
	MonthlyGroundRent          float64                           `json:"monthlyGroundRent"`
	Exchange                   bool                              `json:"exchange"`
	ExchangePercentual         float64                           `json:"exchangePercentual"`
	ExchangePlaces             []ListingExchangePlaceResponse    `json:"exchangePlaces,omitempty"`
	Installment                *CatalogItemResponse              `json:"installment,omitempty"`
	Financing                  bool                              `json:"financing"`
	FinancingBlockers          []ListingFinancingBlockerResponse `json:"financingBlockers,omitempty"`
	Guarantees                 []ListingGuaranteeResponse        `json:"guarantees,omitempty"`
	Visit                      *CatalogItemResponse              `json:"visit,omitempty"`
	TenantName                 string                            `json:"tenantName"`
	TenantEmail                string                            `json:"tenantEmail"`
	TenantPhone                string                            `json:"tenantPhone"`
	Accompanying               *CatalogItemResponse              `json:"accompanying,omitempty"`
	PhotoSessionID             *uint64                           `json:"photoSessionId,omitempty"`
	Deleted                    bool                              `json:"deleted"`
	PerformanceMetrics         ListingPerformanceMetricsResponse `json:"performanceMetrics"`
	FavoritesCount             int64                             `json:"favoritesCount"`
	IsFavorite                 bool                              `json:"isFavorite"`
	CompletionForecast         string                            `json:"completionForecast,omitempty" example:"2026-06"`
	LandBlock                  string                            `json:"landBlock,omitempty" example:"A"`
	LandLot                    string                            `json:"landLot,omitempty" example:"15"`
	LandFront                  float64                           `json:"landFront,omitempty" example:"12.5"`
	LandSide                   float64                           `json:"landSide,omitempty" example:"30.0"`
	LandBack                   float64                           `json:"landBack,omitempty" example:"12.5"`
	LandTerrainType            *CatalogItemResponse              `json:"landTerrainType,omitempty"`
	HasKmz                     bool                              `json:"hasKmz,omitempty"`
	KmzFile                    string                            `json:"kmzFile,omitempty" example:"https://storage.exemplo.com/terrenos/lote15.kmz"`
	BuildingFloors             int16                             `json:"buildingFloors,omitempty" example:"8"`
	UnitTower                  string                            `json:"unitTower,omitempty" example:"Torre B"`
	UnitFloor                  int16                             `json:"unitFloor,omitempty" example:"5"`
	UnitNumber                 string                            `json:"unitNumber,omitempty" example:"502"`
	WarehouseManufacturingArea float64                           `json:"warehouseManufacturingArea,omitempty" example:"850.5"`
	WarehouseSector            *CatalogItemResponse              `json:"warehouseSector,omitempty"`
	WarehouseHasPrimaryCabin   bool                              `json:"warehouseHasPrimaryCabin,omitempty"`
	WarehouseCabinKva          float64                           `json:"warehouseCabinKva,omitempty" example:"150.0"`
	WarehouseGroundFloor       float64                           `json:"warehouseGroundFloor,omitempty" example:"4.2"`
	WarehouseFloorResistance   float64                           `json:"warehouseFloorResistance,omitempty" example:"2500.0"`
	WarehouseZoning            string                            `json:"warehouseZoning,omitempty" example:"ZI-2"`
	WarehouseHasOfficeArea     bool                              `json:"warehouseHasOfficeArea,omitempty"`
	WarehouseOfficeArea        float64                           `json:"warehouseOfficeArea,omitempty" example:"120.0"`
	WarehouseAdditionalFloors  []WarehouseAdditionalFloorDTO     `json:"warehouseAdditionalFloors,omitempty"`
	StoreHasMezzanine          bool                              `json:"storeHasMezzanine,omitempty"`
	StoreMezzanineArea         float64                           `json:"storeMezzanineArea,omitempty" example:"45.0"`
}

// WarehouseAdditionalFloorDTO represents additional floors in warehouses beyond ground floor.
//
// @Description Additional floor information for warehouses (mezanino, second floor, etc.)
// @property floorName string true "Name/identifier of the floor" example("Mezanino")
// @property floorOrder integer true "Ordering position (1=first above ground, 2=second, etc.)" example(1)
// @property floorHeight number true "Ceiling height in meters" example(3.5)
type WarehouseAdditionalFloorDTO struct {
	FloorName   string  `json:"floorName" example:"Mezanino" description:"Name/identifier of the floor"`
	FloorOrder  int     `json:"floorOrder" example:"1" description:"Ordering position (1=first above ground, 2=second, etc.)"`
	FloorHeight float64 `json:"floorHeight" example:"3.5" description:"Ceiling height in meters"`
}

// UpdateListingRequest represents request for updating a listing.
// @Description Request payload for updating draft listing. Omitted fields remain unchanged; present fields (including null) overwrite stored values.
// @Description Property-specific required fields (validated on promote): Casa em Construção requires completionForecast; All Terrenos require landBlock; Terreno Comercial/Residencial require landLot, landTerrainType, hasKmz (kmzFile required if hasKmz=true); Prédio requires buildingFloors; Apartamento/Sala/Laje require unitTower, unitFloor, unitNumber; Galpão requires warehouseManufacturingArea, warehouseSector, warehouseHasPrimaryCabin (warehouseCabinKva required if true), warehouseGroundFloor, warehouseFloorResistance, warehouseZoning, warehouseHasOfficeArea (warehouseOfficeArea required if true); Loja requires storeHasMezzanine (storeMezzanineArea required if true).
//
// Exemplo completo de payload (todos os campos preenchidos):
//
//	{
//	  "listingIdentityId": 1024,
//	  "listingVersionId": 5001,
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
//	  "title": "Apartamento 3 dormitorios com piscina",
//	  "accompanying": "assistant",
//	  "completionForecast": "2026-06",
//	  "landBlock": "A",
//	  "landLot": "15",
//	  "landFront": 12.5,
//	  "landSide": 30.0,
//	  "landBack": 12.5,
//	  "landTerrainType": "plano",
//	  "hasKmz": true,
//	  "kmzFile": "https://storage.exemplo.com/terrenos/lote15.kmz",
//	  "buildingFloors": 8,
//	  "unitTower": "Torre B",
//	  "unitFloor": 5,
//	  "unitNumber": "502",
//	  "warehouseManufacturingArea": 850.5,
//	  "warehouseSector": "industrial",
//	  "warehouseHasPrimaryCabin": true,
//	  "warehouseCabinKva": 150.0,
//	  "warehouseGroundFloor": 4.2,
//	  "warehouseFloorResistance": 2500.0,
//	  "warehouseZoning": "ZI-2",
//	  "warehouseHasOfficeArea": true,
//	  "warehouseOfficeArea": 120.0,
//	  "warehouseAdditionalFloors": [
//	    {"floorName": "Mezanino", "floorOrder": 1, "floorHeight": 3.5},
//	    {"floorName": "Segundo Piso", "floorOrder": 2, "floorHeight": 3.2}
//	  ],
//	  "storeHasMezzanine": true,
//	  "storeMezzanineArea": 45.0
//	}
//
// @property listingIdentityId integer true "Listing identity ID" example(1024)
// @property listingVersionId integer true "Listing version ID" example(5001)
// @property owner string false "Owner type: myself, third_party" example(myself)
// @property features array false "List of features with quantities"
// @property landSize number false "Land size in square meters" example(423.5)
// @property corner boolean false "Whether property is on a corner" example(true)
// @property nonBuildable number false "Non-buildable area in square meters" example(12.75)
// @property buildable number false "Buildable area in square meters" example(410.75)
// @property delivered string false "Delivery state: furnished, unfurnished, semi_furnished" example(furnished)
// @property whoLives string false "Current occupancy: owner, tenant, vacant" example(tenant)
// @property title string false "Listing title" example(Apartamento 3 dormitorios com piscina)
// @property description string false "Listing description" example(Apartamento amplo com vista panoramica)
// @property transaction string false "Transaction type: sale, rent, both" example(sale)
// @property sellNet number false "Sale price" example(1200000)
// @property rentNet number false "Monthly rent price" example(8500)
// @property complex string false "Name of the condominium/complex" example("Residencial Jardins")
// @property condominium number false "Monthly condominium fee" example(1200.5)
// @property annualTax number false "Annual IPTU (property tax). Mutually exclusive with monthlyTax" example(3400.75)
// @property monthlyTax number false "Monthly IPTU (property tax). Mutually exclusive with annualTax" example(283.40)
// @property annualGroundRent number false "Annual Laudêmio (ground rent). Mutually exclusive with monthlyGroundRent" example(1800)
// @property monthlyGroundRent number false "Monthly Laudêmio (ground rent). Mutually exclusive with annualGroundRent" example(150)
// @property exchange boolean false "Whether owner accepts exchange" example(true)
// @property exchangePercentual number false "Exchange percentage accepted" example(50)
// @property exchangePlaces array false "List of acceptable exchange locations"
// @property installment string false "Installment type: short_term, long_term, none" example(short_term)
// @property financing boolean false "Whether financing is available" example(true)
// @property financingBlockers array false "List of financing blocker codes"
// @property guarantees array false "List of accepted guarantees with priority"
// @property visit string false "Visit policy: owner, client, flexible" example(client)
// @property tenantName string false "Current tenant name" example(Joao da Silva)
// @property tenantEmail string false "Current tenant email" example(joao.silva@example.com)
// @property tenantPhone string false "Current tenant phone in E.164 format" example(+5511912345678)
// @property accompanying string false "Accompanying requirement: broker, assistant, owner, none" example(assistant)
// @property completionForecast string false "Completion forecast for properties under construction. Accepts: YYYY-MM-DD (2026-01-20), YYYY-MM (2026-06), or RFC3339 timestamp (2026-01-20T00:00:00Z). Internally normalized to YYYY-MM-DD for storage." example(2026-06)
// @property landBlock string false "Block identifier for land properties" example(A)
// @property landLot string false "Lot number for land properties" example(15)
// @property landFront number false "Front dimension in meters for land properties" example(12.5)
// @property landSide number false "Side dimension in meters for land properties" example(30.0)
// @property landBack number false "Back dimension in meters for land properties" example(12.5)
// @property landTerrainType string false "Terrain type: plano, aclive, declive, irregular, misto" example(plano)
// @property hasKmz boolean false "Indicates if KMZ file is available for land properties" example(true)
// @property kmzFile string false "URL to KMZ file for land properties" example(https://storage.exemplo.com/terrenos/lote15.kmz)
// @property buildingFloors integer false "Total number of floors in building" example(8)
// @property unitTower string false "Tower identifier for apartment/commercial units" example(Torre B)
// @property unitFloor integer false "Floor number where unit is located" example(5)
// @property unitNumber string false "Unit number/identifier" example(502)
// @property warehouseManufacturingArea number false "Manufacturing/production area in square meters for warehouses" example(850.5)
// @property warehouseSector string false "Warehouse sector: industrial, logistico, comercial" example(industrial)
// @property warehouseHasPrimaryCabin boolean false "Indicates if warehouse has primary electrical cabin" example(true)
// @property warehouseCabinKva number false "Primary cabin power in KVA" example(150.0)
// @property warehouseGroundFloor number false "Ground floor ceiling height in meters" example(4.2)
// @property warehouseFloorResistance number false "Floor resistance in kg/m²" example(2500.0)
// @property warehouseZoning string false "Zoning classification for warehouse" example(ZI-2)
// @property warehouseHasOfficeArea boolean false "Indicates if warehouse has office area" example(true)
// @property warehouseOfficeArea number false "Office area in square meters" example(120.0)
// @property warehouseAdditionalFloors array false "Additional floors beyond ground floor in warehouses"
// @property storeHasMezzanine boolean false "Indicates if store has mezzanine" example(true)
// @property storeMezzanineArea number false "Mezzanine area in square meters for stores" example(45.0)
type UpdateListingRequest struct {
	ListingIdentityID          coreutils.Optional[int64]                               `json:"listingIdentityId" binding:"required" example:"1024"`
	ListingVersionID           coreutils.Optional[int64]                               `json:"listingVersionId" binding:"required" example:"5001"`
	Owner                      coreutils.Optional[string]                              `json:"owner" example:"myself"`
	Features                   coreutils.Optional[[]UpdateListingFeatureRequest]       `json:"features"`
	LandSize                   coreutils.Optional[float64]                             `json:"landSize" example:"423.5"`
	Corner                     coreutils.Optional[bool]                                `json:"corner" example:"true"`
	NonBuildable               coreutils.Optional[float64]                             `json:"nonBuildable" example:"12.75"`
	Buildable                  coreutils.Optional[float64]                             `json:"buildable" example:"410.75"`
	Delivered                  coreutils.Optional[string]                              `json:"delivered" example:"furnished"`
	WhoLives                   coreutils.Optional[string]                              `json:"whoLives" example:"tenant"`
	Title                      coreutils.Optional[string]                              `json:"title" example:"Apartamento 3 dormitorios com piscina"`
	Description                coreutils.Optional[string]                              `json:"description" example:"Apartamento amplo com vista panoramica"`
	Complex                    coreutils.Optional[string]                              `json:"complex" example:"Residencial Jardins" description:"Name of the condominium/complex"`
	Transaction                coreutils.Optional[string]                              `json:"transaction" example:"sale"`
	SellNet                    coreutils.Optional[float64]                             `json:"sellNet" example:"1200000"`
	RentNet                    coreutils.Optional[float64]                             `json:"rentNet" example:"8500"`
	Condominium                coreutils.Optional[float64]                             `json:"condominium" example:"1200.5"`
	AnnualTax                  coreutils.Optional[float64]                             `json:"annualTax" example:"3400.75" description:"Annual IPTU (property tax). Mutually exclusive with monthlyTax."`
	MonthlyTax                 coreutils.Optional[float64]                             `json:"monthlyTax" example:"283.40" description:"Monthly IPTU (property tax). Mutually exclusive with annualTax."`
	AnnualGroundRent           coreutils.Optional[float64]                             `json:"annualGroundRent" example:"1800" description:"Annual Laudêmio (ground rent). Mutually exclusive with monthlyGroundRent."`
	MonthlyGroundRent          coreutils.Optional[float64]                             `json:"monthlyGroundRent" example:"150" description:"Monthly Laudêmio (ground rent). Mutually exclusive with annualGroundRent."`
	Exchange                   coreutils.Optional[bool]                                `json:"exchange" example:"true"`
	ExchangePercentual         coreutils.Optional[float64]                             `json:"exchangePercentual" example:"50"`
	ExchangePlaces             coreutils.Optional[[]UpdateListingExchangePlaceRequest] `json:"exchangePlaces"`
	Installment                coreutils.Optional[string]                              `json:"installment" example:"short_term"`
	Financing                  coreutils.Optional[bool]                                `json:"financing" example:"true"`
	FinancingBlockers          coreutils.Optional[[]string]                            `json:"financingBlockers"`
	Guarantees                 coreutils.Optional[[]UpdateListingGuaranteeRequest]     `json:"guarantees"`
	Visit                      coreutils.Optional[string]                              `json:"visit" example:"client"`
	TenantName                 coreutils.Optional[string]                              `json:"tenantName" example:"Joao da Silva"`
	TenantEmail                coreutils.Optional[string]                              `json:"tenantEmail" example:"joao.silva@example.com"`
	TenantPhone                coreutils.Optional[string]                              `json:"tenantPhone" example:"+5511912345678" description:"Tenant phone number in E.164 format (e.g., +5511912345678)."`
	Accompanying               coreutils.Optional[string]                              `json:"accompanying" example:"assistant"`
	CompletionForecast         coreutils.Optional[string]                              `json:"completionForecast" example:"2026-06" description:"Completion forecast for properties under construction. Accepts: YYYY-MM-DD, YYYY-MM, or RFC3339 timestamp. Internally normalized to YYYY-MM-DD for storage."`
	LandBlock                  coreutils.Optional[string]                              `json:"landBlock" example:"A" description:"Block identifier for land properties"`
	LandLot                    coreutils.Optional[string]                              `json:"landLot" example:"15" description:"Lot number for land properties"`
	LandFront                  coreutils.Optional[float64]                             `json:"landFront" example:"12.5" description:"Front dimension in meters for land properties"`
	LandSide                   coreutils.Optional[float64]                             `json:"landSide" example:"30.0" description:"Side dimension in meters for land properties"`
	LandBack                   coreutils.Optional[float64]                             `json:"landBack" example:"12.5" description:"Back dimension in meters for land properties"`
	LandTerrainType            coreutils.Optional[string]                              `json:"landTerrainType" example:"plano" description:"Terrain type: plano, aclive, declive, irregular, misto"`
	HasKmz                     coreutils.Optional[bool]                                `json:"hasKmz" example:"true" description:"Indicates if KMZ file is available for land properties"`
	KmzFile                    coreutils.Optional[string]                              `json:"kmzFile" example:"https://storage.exemplo.com/terrenos/lote15.kmz" description:"URL to KMZ file for land properties"`
	BuildingFloors             coreutils.Optional[int16]                               `json:"buildingFloors" example:"8" description:"Total number of floors in building"`
	UnitTower                  coreutils.Optional[string]                              `json:"unitTower" example:"Torre B" description:"Tower identifier for apartment/commercial units"`
	UnitFloor                  coreutils.Optional[int16]                               `json:"unitFloor" example:"5" description:"Floor number where unit is located"`
	UnitNumber                 coreutils.Optional[string]                              `json:"unitNumber" example:"502" description:"Unit number/identifier"`
	WarehouseManufacturingArea coreutils.Optional[float64]                             `json:"warehouseManufacturingArea" example:"850.5" description:"Manufacturing/production area in square meters for warehouses"`
	WarehouseSector            coreutils.Optional[string]                              `json:"warehouseSector" example:"industrial" description:"Warehouse sector: industrial, logistico, comercial"`
	WarehouseHasPrimaryCabin   coreutils.Optional[bool]                                `json:"warehouseHasPrimaryCabin" example:"true" description:"Indicates if warehouse has primary electrical cabin"`
	WarehouseCabinKva          coreutils.Optional[float64]                             `json:"warehouseCabinKva" example:"150.0" description:"Primary cabin power in KVA"`
	WarehouseGroundFloor       coreutils.Optional[float64]                             `json:"warehouseGroundFloor" example:"4.2" description:"Ground floor ceiling height in meters"`
	WarehouseFloorResistance   coreutils.Optional[float64]                             `json:"warehouseFloorResistance" example:"2500.0" description:"Floor resistance in kg/m²"`
	WarehouseZoning            coreutils.Optional[string]                              `json:"warehouseZoning" example:"ZI-2" description:"Zoning classification for warehouse"`
	WarehouseHasOfficeArea     coreutils.Optional[bool]                                `json:"warehouseHasOfficeArea" example:"true" description:"Indicates if warehouse has office area"`
	WarehouseOfficeArea        coreutils.Optional[float64]                             `json:"warehouseOfficeArea" example:"120.0" description:"Office area in square meters"`
	WarehouseAdditionalFloors  coreutils.Optional[[]WarehouseAdditionalFloorDTO]       `json:"warehouseAdditionalFloors" description:"Additional floors beyond ground floor in warehouses"`
	StoreHasMezzanine          coreutils.Optional[bool]                                `json:"storeHasMezzanine" example:"true" description:"Indicates if store has mezzanine"`
	StoreMezzanineArea         coreutils.Optional[float64]                             `json:"storeMezzanineArea" example:"45.0" description:"Mezzanine area in square meters for stores"`
}

// UpdateListingRequestSwagger is used ONLY for Swagger documentation since swag doesn't support Optional[T] generics.
// The actual handler uses UpdateListingRequest with Optional fields.
// This struct mirrors UpdateListingRequest but with concrete types to generate proper Swagger docs.
type UpdateListingRequestSwagger struct {
	ListingIdentityID          int64                                `json:"listingIdentityId" binding:"required" example:"1024"`
	ListingVersionID           int64                                `json:"listingVersionId" binding:"required" example:"5001"`
	Owner                      *string                              `json:"owner,omitempty" example:"myself"`
	Features                   *[]UpdateListingFeatureRequest       `json:"features,omitempty"`
	LandSize                   *float64                             `json:"landSize,omitempty" example:"423.5"`
	Corner                     *bool                                `json:"corner,omitempty" example:"true"`
	NonBuildable               *float64                             `json:"nonBuildable,omitempty" example:"12.75"`
	Buildable                  *float64                             `json:"buildable,omitempty" example:"410.75"`
	Delivered                  *string                              `json:"delivered,omitempty" example:"furnished"`
	WhoLives                   *string                              `json:"whoLives,omitempty" example:"tenant"`
	Title                      *string                              `json:"title,omitempty" example:"Apartamento 3 dormitorios com piscina"`
	Description                *string                              `json:"description,omitempty" example:"Apartamento amplo com vista panoramica"`
	Complex                    *string                              `json:"complex,omitempty" example:"Residencial Jardins"`
	Transaction                *string                              `json:"transaction,omitempty" example:"sale"`
	SellNet                    *float64                             `json:"sellNet,omitempty" example:"1200000"`
	RentNet                    *float64                             `json:"rentNet,omitempty" example:"8500"`
	Condominium                *float64                             `json:"condominium,omitempty" example:"1200.5"`
	AnnualTax                  *float64                             `json:"annualTax,omitempty" example:"3400.75"`
	MonthlyTax                 *float64                             `json:"monthlyTax,omitempty" example:"283.40"`
	AnnualGroundRent           *float64                             `json:"annualGroundRent,omitempty" example:"1800"`
	MonthlyGroundRent          *float64                             `json:"monthlyGroundRent,omitempty" example:"150"`
	Exchange                   *bool                                `json:"exchange,omitempty" example:"true"`
	ExchangePercentual         *float64                             `json:"exchangePercentual,omitempty" example:"50"`
	ExchangePlaces             *[]UpdateListingExchangePlaceRequest `json:"exchangePlaces,omitempty"`
	Installment                *string                              `json:"installment,omitempty" example:"short_term"`
	Financing                  *bool                                `json:"financing,omitempty" example:"true"`
	FinancingBlockers          *[]string                            `json:"financingBlockers,omitempty"`
	Guarantees                 *[]UpdateListingGuaranteeRequest     `json:"guarantees,omitempty"`
	Visit                      *string                              `json:"visit,omitempty" example:"client"`
	TenantName                 *string                              `json:"tenantName,omitempty" example:"Joao da Silva"`
	TenantEmail                *string                              `json:"tenantEmail,omitempty" example:"joao.silva@example.com"`
	TenantPhone                *string                              `json:"tenantPhone,omitempty" example:"+5511912345678"`
	Accompanying               *string                              `json:"accompanying,omitempty" example:"assistant"`
	CompletionForecast         *string                              `json:"completionForecast,omitempty" example:"2026-06" description:"Completion forecast. Accepts: YYYY-MM-DD, YYYY-MM, or RFC3339. Normalized to YYYY-MM-DD."`
	LandBlock                  *string                              `json:"landBlock,omitempty" example:"A"`
	LandLot                    *string                              `json:"landLot,omitempty" example:"15"`
	LandFront                  *float64                             `json:"landFront,omitempty" example:"12.5"`
	LandSide                   *float64                             `json:"landSide,omitempty" example:"30.0"`
	LandBack                   *float64                             `json:"landBack,omitempty" example:"12.5"`
	LandTerrainType            *string                              `json:"landTerrainType,omitempty" example:"plano"`
	HasKmz                     *bool                                `json:"hasKmz,omitempty" example:"true"`
	KmzFile                    *string                              `json:"kmzFile,omitempty" example:"https://storage.exemplo.com/terrenos/lote15.kmz"`
	BuildingFloors             *int16                               `json:"buildingFloors,omitempty" example:"8"`
	UnitTower                  *string                              `json:"unitTower,omitempty" example:"Torre B"`
	UnitFloor                  *int16                               `json:"unitFloor,omitempty" example:"5"`
	UnitNumber                 *string                              `json:"unitNumber,omitempty" example:"502"`
	WarehouseManufacturingArea *float64                             `json:"warehouseManufacturingArea,omitempty" example:"850.5"`
	WarehouseSector            *string                              `json:"warehouseSector,omitempty" example:"industrial"`
	WarehouseHasPrimaryCabin   *bool                                `json:"warehouseHasPrimaryCabin,omitempty" example:"true"`
	WarehouseCabinKva          *float64                             `json:"warehouseCabinKva,omitempty" example:"150.0"`
	WarehouseGroundFloor       *float64                             `json:"warehouseGroundFloor,omitempty" example:"4.2"`
	WarehouseFloorResistance   *float64                             `json:"warehouseFloorResistance,omitempty" example:"2500.0"`
	WarehouseZoning            *string                              `json:"warehouseZoning,omitempty" example:"ZI-2"`
	WarehouseHasOfficeArea     *bool                                `json:"warehouseHasOfficeArea,omitempty" example:"true"`
	WarehouseOfficeArea        *float64                             `json:"warehouseOfficeArea,omitempty" example:"120.0"`
	WarehouseAdditionalFloors  *[]WarehouseAdditionalFloorDTO       `json:"warehouseAdditionalFloors,omitempty"`
	StoreHasMezzanine          *bool                                `json:"storeHasMezzanine,omitempty" example:"true"`
	StoreMezzanineArea         *float64                             `json:"storeMezzanineArea,omitempty" example:"45.0"`
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

// GetOptionsResponse represents response for getting options
type GetOptionsResponse struct {
	PropertyTypes []PropertyTypeOption `json:"propertyTypes"`
	ComplexName   string               `json:"complexName,omitempty"`
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
	ID                int64                        `json:"id"`
	ListingIdentityID int64                        `json:"listingIdentityId"`
	ListingUUID       string                       `json:"listingUuid"`
	ActiveVersionID   int64                        `json:"activeVersionId"`
	DraftVersionID    *int64                       `json:"draftVersionId,omitempty"`
	Version           uint8                        `json:"version"`
	Title             string                       `json:"title"`
	Description       string                       `json:"description"`
	Price             float64                      `json:"price"`
	Status            string                       `json:"status"`
	PropertyType      *ListingPropertyTypeResponse `json:"propertyType,omitempty"`
	ZipCode           string                       `json:"zipCode"`
	Number            string                       `json:"number"`
	Complex           string                       `json:"complex,omitempty"`
	UserID            int64                        `json:"userId"`
	ComplexID         string                       `json:"complexId,omitempty"`
	FavoritesCount    int64                        `json:"favoritesCount"`
	IsFavorite        bool                         `json:"isFavorite"`
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
	ListingIdentityID int64  `json:"listingIdentityId" binding:"required" example:"1024"`
	PhotoSessionID    uint64 `json:"photoSessionId" binding:"required" example:"3003"`
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

// ====================================================================================================
// Media Processing DTOs
// ====================================================================================================

// RequestUploadURLsRequest represents the payload used to request signed URLs for raw media uploads
//
// This DTO is used to initialize a media upload batch for a listing.
// The service validates the manifest, generates signed S3 URLs, and returns upload instructions.
type RequestUploadURLsRequest struct {
	// ListingIdentityID identifies the listing receiving the batch
	// Must correspond to a listing in PENDING_PHOTO_PROCESSING status
	// Example: 123
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"123"`

	// Files enumerates every asset to upload, preserving carousel order and metadata
	// Minimum 1 file, maximum defined by environment configuration (default 60)
	// Each file must have unique clientId and sequence
	Files []RequestUploadFileRequest `json:"files" binding:"required,min=1,dive"`
}

// RequestUploadFileRequest describes a single asset to be uploaded by the client
//
// Each file in the manifest must provide complete metadata for validation and storage organization.
type RequestUploadFileRequest struct {
	// AssetType categorizes the media for processing and display
	// Allowed values: PHOTO_VERTICAL, PHOTO_HORIZONTAL, VIDEO_VERTICAL, VIDEO_HORIZONTAL,
	//                 THUMBNAIL, ZIP, PROJECT_DOC, PROJECT_RENDER
	// Example: "PHOTO_VERTICAL"
	AssetType string `json:"assetType" binding:"required" example:"PHOTO_VERTICAL"`

	// Orientation specifies the visual orientation of the asset
	// Allowed values: PORTRAIT, LANDSCAPE
	// Example: "PORTRAIT"
	Orientation string `json:"orientation" binding:"required" example:"PORTRAIT"`

	// Filename is the original filename from the client device
	// Used for download metadata and logging
	// Example: "IMG_20231116_140530.jpg"
	Filename string `json:"filename" binding:"required" example:"IMG_20231116_140530.jpg"`

	// ContentType is the MIME type of the file
	// Must be in the allowed content types list (configured in env.yaml)
	// Common values: image/jpeg, image/png, image/heic, video/mp4, video/quicktime
	// Example: "image/jpeg"
	ContentType string `json:"contentType" binding:"required" example:"image/jpeg"`

	// Bytes is the file size in bytes
	// Must be positive and not exceed the maximum file size limit
	// Example: 2457600 (2.4 MB)
	Bytes int64 `json:"bytes" binding:"required,min=1" example:"2457600"`

	// Checksum is the SHA-256 hash of the file content (hex-encoded)
	// Used to verify upload integrity via S3 HeadObject
	// Example: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	Checksum string `json:"checksum" binding:"required" example:"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"`

	// Title is an optional user-provided caption or description
	// Displayed in the listing carousel
	// Example: "Vista frontal do imóvel"
	Title string `json:"title,omitempty" example:"Vista frontal do imóvel"`

	// Sequence determines the display order in the listing carousel
	// Must be positive and unique within the batch
	// Example: 1
	Sequence uint8 `json:"sequence" binding:"required,min=1" example:"1"`

	// Metadata holds additional key-value pairs for processing hints
	// Optional field for future extensibility
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RequestUploadURLsResponse returns signed URLs ready to be used by the uploader
//
// The client must use these signed URLs to upload files to S3.
type RequestUploadURLsResponse struct {
	// ListingIdentityID confirms the listing receiving the uploads
	ListingIdentityID uint64 `json:"listingIdentityId" example:"123"`

	// UploadURLTTLSeconds indicates how long the signed URLs remain valid
	// Typically 900 seconds (15 minutes)
	UploadURLTTLSeconds int `json:"uploadUrlTtlSeconds" example:"900"`

	// Files contains upload instructions for each asset in the manifest
	Files []RequestUploadInstructionResponse `json:"files"`
}

// UploadInstructionResponse carries the information required to perform a PUT upload to S3
//
// The client must execute an HTTP PUT request to UploadURL with the specified headers.
type RequestUploadInstructionResponse struct {
	// AssetType matches the assetType from the request
	AssetType string `json:"assetType" example:"PHOTO_VERTICAL"`

	// Sequence matches the sequence from the request
	Sequence uint8 `json:"sequence" example:"1"`

	// UploadURL is the pre-signed S3 URL for uploading the file
	// Valid for the duration specified in UploadURLTTLSeconds
	// Example: "https://s3.amazonaws.com/bucket/key?X-Amz-Algorithm=..."
	UploadURL string `json:"uploadUrl" example:"https://s3.amazonaws.com/bucket/key?X-Amz-Algorithm=..."`

	// Method is the HTTP method to use for the upload (always PUT)
	Method string `json:"method" example:"PUT"`

	// Headers contains required HTTP headers for the upload request
	// Typically includes Content-Type and x-amz-checksum-sha256
	Headers map[string]string `json:"headers"`

	// ObjectKey is the S3 object key where the file will be stored
	// Format: /{listingIdentityId}/raw/{assetType}/{YYYY-MM-DD}/{filename}
	// Example: "123/raw/PHOTO_VERTICAL/2025-11-16/IMG_20231116_140530.jpg"
	ObjectKey string `json:"objectKey" example:"123/raw/PHOTO_VERTICAL/2025-11-16/IMG_20231116_140530.jpg"`

	// Title matches the title from the request
	Title string `json:"title,omitempty" example:"Vista frontal do imóvel"`
}

// CompleteUploadBatchRequest confirms that every asset in the manifest has been uploaded successfully
//
// The client must call this endpoint after all files have been uploaded to S3.
// The service validates each upload via S3 HeadObject, then enqueues the batch for async processing.
type CompleteUploadBatchRequest struct {
	// ListingIdentityID must match the listingIdentityId from CreateUploadBatch
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"123"`

	// BatchID must match the batchId from CreateUploadBatch
	BatchID uint64 `json:"batchId" binding:"required,min=1" example:"456"`

	// Files lists all successfully uploaded assets with S3 metadata
	// Must include all files from the original manifest
	Files []CompletedUploadFileRequest `json:"files" binding:"required,min=1,dive"`
}

// CompletedUploadFileRequest links the logical asset to the physical object persisted in S3
//
// The client provides confirmation that the upload succeeded and includes S3 ETag for verification.
type CompletedUploadFileRequest struct {
	// ClientID must match a clientId from the original CreateUploadBatch request
	ClientID string `json:"clientId" binding:"required" example:"photo-001"`

	// ObjectKey must match the objectKey returned in CreateUploadBatch
	ObjectKey string `json:"objectKey" binding:"required" example:"123/raw/PHOTO_VERTICAL/2025-11-16/IMG_20231116_140530.jpg"`

	// Bytes is the confirmed file size after upload (should match original)
	Bytes int64 `json:"bytes" binding:"required,min=1" example:"2457600"`

	// Checksum is the SHA-256 hash used for upload verification
	Checksum string `json:"checksum" binding:"required" example:"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"`

	// ETag is the S3 ETag returned after successful upload
	// Optional but recommended for additional verification
	ETag string `json:"etag,omitempty" example:"\"5d41402abc4b2a76b9719d911017c592\""`
}

// CompleteUploadBatchResponse exposes the async job metadata published to the processing queue
//
// The batch transitions to RECEIVED status and a processing job is enqueued.
// The client should poll GetBatchStatus to monitor progress.
type CompleteUploadBatchResponse struct {
	// ListingIdentityID confirms the listing
	ListingIdentityID uint64 `json:"listingIdentityId" example:"123"`

	// BatchID confirms the batch
	BatchID uint64 `json:"batchId" example:"456"`

	// JobID is the unique identifier for the async processing job
	// Used for monitoring and troubleshooting
	JobID uint64 `json:"jobId" example:"789"`

	// Status indicates the current batch status (typically RECEIVED)
	Status string `json:"status" example:"RECEIVED"`

	// EstimatedDurationSeconds provides a rough estimate of processing time
	// Actual duration depends on file count, sizes, and pipeline load
	EstimatedDurationSeconds int `json:"estimatedDurationSeconds" example:"300"`
}

// GetBatchStatusRequest retrieves detailed status for a specific batch under a listing identity
//
// Used by the frontend to poll batch progress during upload and processing.
type GetBatchStatusRequest struct {
	// ListingIdentityID identifies the listing owning the batch
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"123"`

	// BatchID identifies the specific batch to query
	BatchID uint64 `json:"batchId" binding:"required,min=1" example:"456"`
}

// GetBatchStatusResponse aggregates batch status and asset metadata for UI polling
//
// Contains current batch status and detailed information about each asset.
type GetBatchStatusResponse struct {
	// ListingIdentityID confirms the listing
	ListingIdentityID uint64 `json:"listingIdentityId" example:"123"`

	// BatchID confirms the batch
	BatchID uint64 `json:"batchId" example:"456"`

	// Status indicates the current batch status
	// Possible values: PENDING_UPLOAD, RECEIVED, PROCESSING, READY, FAILED
	Status string `json:"status" example:"PROCESSING"`

	// StatusMessage provides human-readable status information
	// Includes error details for FAILED status
	StatusMessage string `json:"statusMessage" example:"Processing assets"`

	// Assets lists all assets in the batch with their processing status
	Assets []BatchAssetStatusResponse `json:"assets"`
}

// BatchAssetStatusResponse mirrors the information frontend needs to display upload and processing progress
//
// Each asset shows upload status and processed output availability.
type BatchAssetStatusResponse struct {
	// ClientID matches the original clientId for correlation
	ClientID string `json:"clientId" example:"photo-001"`

	// Title matches the original title
	Title string `json:"title,omitempty" example:"Vista frontal do imóvel"`

	// AssetType matches the original assetType
	AssetType string `json:"assetType" example:"PHOTO_VERTICAL"`

	// Sequence matches the original sequence
	Sequence uint8 `json:"sequence" example:"1"`

	// RawObjectKey is the S3 key where the raw file is stored
	// Empty until upload is confirmed
	RawObjectKey string `json:"rawObjectKey,omitempty" example:"123/raw/PHOTO_VERTICAL/2025-11-16/IMG_20231116_140530.jpg"`

	// ProcessedKey is the S3 key where the processed file is stored
	// Empty until processing completes successfully
	ProcessedKey string `json:"processedKey,omitempty" example:"123/processed/PHOTO_VERTICAL/2025-11-16/IMG_20231116_140530.jpg"`

	// ThumbnailKey is the S3 key where the thumbnail is stored
	// Empty until processing completes successfully
	ThumbnailKey string `json:"thumbnailKey,omitempty" example:"123/processed/THUMBNAIL/2025-11-16/IMG_20231116_140530_thumb.jpg"`

	// Metadata contains additional asset information
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListDownloadURLsRequest requests signed GET URLs for processed assets
//
// Used to retrieve download links for assets that have completed processing.
type ListDownloadURLsRequest struct {
	// ListingIdentityID identifies the listing owning the assets
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"123"`

	// AssetTypes filters the results by asset type
	AssetTypes []string `json:"assetTypes,omitempty" example:"PHOTO_VERTICAL,VIDEO_HORIZONTAL"`
}

// ListDownloadURLsResponse returns signed URLs for processed assets within a batch
//
// Contains download links valid for the duration specified in TTLSeconds.
type ListDownloadURLsResponse struct {
	// ListingIdentityID confirms the listing
	ListingIdentityID uint64 `json:"listingIdentityId" example:"123"`

	// Downloads contains signed URLs for each processed asset
	Downloads []DownloadEntryResponse `json:"downloads"`
}

// DownloadEntryResponse encapsulates individual signed URLs derived from processed assets
//
// Each entry provides download links for the processed file and optional preview/thumbnail.
type DownloadEntryResponse struct {

	// AssetType matches the original assetType
	AssetType string `json:"assetType" example:"PHOTO_VERTICAL"`

	// Sequence matches the original sequence
	Sequence uint8 `json:"sequence" example:"1"`

	// Status indicates the processing status
	Status string `json:"status" example:"PROCESSED"`

	// Title matches the original title
	Title string `json:"title,omitempty" example:"Vista frontal do imóvel"`

	// URL is the signed GET URL for the processed file
	URL string `json:"url,omitempty" example:"https://s3.amazonaws.com/bucket/key?X-Amz-Algorithm=..."`

	// PreviewURL is the signed GET URL for the thumbnail (if available)
	PreviewURL string `json:"previewUrl,omitempty" example:"https://s3.amazonaws.com/bucket/thumb?X-Amz-Algorithm=..."`

	// Metadata contains additional asset information
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RetryMediaBatchRequest allows the caller to re-enqueue a finished batch
//
// Used to retry processing for batches that failed or need reprocessing.
// Only accepts batches in terminal states (READY or FAILED).
type RetryMediaBatchRequest struct {
	// ListingIdentityID identifies the listing owning the batch
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"123"`

	// BatchID identifies the specific batch to retry
	BatchID uint64 `json:"batchId" binding:"required,min=1" example:"456"`

	// Reason explains why the batch is being retried
	// Used for audit logging and troubleshooting
	// Example: "Pipeline timeout - retrying with increased timeout"
	Reason string `json:"reason" binding:"required,max=500" example:"Pipeline timeout - retrying with increased timeout"`
}

// RetryMediaBatchResponse exposes the identifier of the newly created processing job
//
// The batch transitions back to PROCESSING status and a new job is enqueued.
type RetryMediaBatchResponse struct {
	// ListingIdentityID confirms the listing
	ListingIdentityID uint64 `json:"listingIdentityId" example:"123"`

	// BatchID confirms the batch
	BatchID uint64 `json:"batchId" example:"456"`

	// JobID is the unique identifier for the new async processing job
	JobID uint64 `json:"jobId" example:"790"`

	// Status indicates the current batch status (typically PROCESSING)
	Status string `json:"status" example:"PROCESSING"`
}

// --- List Media (GET) ---

// ListMediaRequest define filtros e paginação para consulta de mídias.
type ListMediaRequest struct {
	ListingIdentityID uint64 `form:"listingIdentityId" binding:"required,min=1"`
	AssetType         string `form:"assetType,omitempty"`
	Sequence          *uint8 `form:"sequence,omitempty"`

	// Paginação e Ordenação
	Page  int    `form:"page,default=1" binding:"min=1"`
	Limit int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort  string `form:"sort,default=sequence" binding:"omitempty,oneof=sequence id"`
	Order string `form:"order,default=asc" binding:"omitempty,oneof=asc desc"`
}

// ListMediaResponse retorna a lista paginada com TODAS as informações do modelo.
type ListMediaResponse struct {
	Data       []MediaAssetResponse    `json:"data"`
	Pagination PaginationResponse      `json:"pagination"`
	ZipBundle  *MediaZipBundleResponse `json:"zipBundle,omitempty"`
}

// ListingMediaApprovalRequest represents owner approval/rejection payload.
type ListingMediaApprovalRequest struct {
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"42"`
	Approve           bool   `json:"approve" example:"true"`
}

// ListingMediaApprovalResponse echoes the new status after the decision.
type ListingMediaApprovalResponse struct {
	ListingIdentityID uint64 `json:"listingIdentityId" example:"42"`
	Decision          string `json:"decision" example:"approved"`
	NewStatus         string `json:"newStatus" example:"PENDING_ADMIN_REVIEW"`
}

// MediaAssetResponse espelha o modelo de domínio completo.
type MediaAssetResponse struct {
	ID                uint64            `json:"id"`
	ListingIdentityID uint64            `json:"listingIdentityId"`
	AssetType         string            `json:"assetType"`
	Sequence          uint8             `json:"sequence"`
	Status            string            `json:"status"`
	Title             string            `json:"title,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	S3KeyRaw          string            `json:"s3KeyRaw,omitempty"`
	S3KeyProcessed    string            `json:"s3KeyProcessed,omitempty"`
}

// MediaZipBundleResponse expõe os metadados do bundle zipado.
type MediaZipBundleResponse struct {
	BundleKey               string `json:"bundleKey"`
	AssetsCount             int    `json:"assetsCount"`
	ZipSizeBytes            int64  `json:"zipSizeBytes"`
	EstimatedExtractedBytes int64  `json:"estimatedExtractedBytes"`
	CompletedAt             string `json:"completedAt,omitempty"`
}

// --- Generate Download URLs (POST) ---

// GenerateDownloadURLsRequest solicita URLs assinadas para itens específicos.
type GenerateDownloadURLsRequest struct {
	ListingIdentityID uint64                `json:"listingIdentityId" binding:"required,min=1"`
	Requests          []DownloadRequestItem `json:"requests" binding:"required,min=1,dive"`
}

// DownloadRequestItem combina a chave do asset com a resolução desejada.
type DownloadRequestItem struct {
	AssetType string `json:"assetType" binding:"required" example:"PHOTO_VERTICAL"`
	Sequence  uint8  `json:"sequence" binding:"required" example:"1"`
	// Resolution options: thumbnail, small, medium, large, original
	Resolution string `json:"resolution" binding:"required,oneof=thumbnail small medium large original" enums:"thumbnail,small,medium,large,original" example:"medium"`
}

// GenerateDownloadURLsResponse retorna as URLs geradas.
type GenerateDownloadURLsResponse struct {
	ListingIdentityID uint64                `json:"listingIdentityId"`
	Urls              []DownloadURLResponse `json:"urls"`
}

type DownloadURLResponse struct {
	AssetType  string `json:"assetType"`
	Sequence   uint8  `json:"sequence"`
	Resolution string `json:"resolution"`
	Url        string `json:"url"`
	ExpiresIn  int    `json:"expiresIn"`
}
