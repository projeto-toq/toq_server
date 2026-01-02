package listingrepository

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

type ListingRepoPortInterface interface {
	// Version-aware operations (Phase 2)
	CreateListingIdentity(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) error
	CreateListingVersion(ctx context.Context, tx *sql.Tx, version listingmodel.ListingVersionInterface) error
	SetListingActiveVersion(ctx context.Context, tx *sql.Tx, identityID int64, versionID int64) error
	GetListingIdentityByID(ctx context.Context, tx *sql.Tx, identityID int64) (ListingIdentityRecord, error)
	GetListingIdentityByUUID(ctx context.Context, tx *sql.Tx, listingUUID string) (ListingIdentityRecord, error)
	GetListingVersionByID(ctx context.Context, tx *sql.Tx, versionID int64) (listingmodel.ListingInterface, error)
	GetDraftVersionByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (listingmodel.ListingInterface, error)
	ListListingVersions(ctx context.Context, tx *sql.Tx, filter ListListingVersionsFilter) ([]ListingVersionSummary, error)

	// Satellite entities
	CreateExchangePlace(ctx context.Context, tx *sql.Tx, place listingmodel.ExchangePlaceInterface) (err error)
	CreateFeature(ctx context.Context, tx *sql.Tx, feature listingmodel.FeatureInterface) (err error)
	CreateGuarantee(ctx context.Context, tx *sql.Tx, guarantee listingmodel.GuaranteeInterface) (err error)
	CreateFinancingBlocker(ctx context.Context, tx *sql.Tx, blocker listingmodel.FinancingBlockerInterface) (err error)
	CreateWarehouseAdditionalFloor(ctx context.Context, tx *sql.Tx, floor listingmodel.WarehouseAdditionalFloorInterface) (err error)

	UpdateExchangePlaces(ctx context.Context, tx *sql.Tx, listingVersionID int64, places []listingmodel.ExchangePlaceInterface) (err error)
	UpdateFeatures(ctx context.Context, tx *sql.Tx, listingVersionID int64, features []listingmodel.FeatureInterface) (err error)
	UpdateGuarantees(ctx context.Context, tx *sql.Tx, listingVersionID int64, guarantees []listingmodel.GuaranteeInterface) (err error)
	UpdateFinancingBlockers(ctx context.Context, tx *sql.Tx, listingVersionID int64, blockers []listingmodel.FinancingBlockerInterface) (err error)
	UpdateWarehouseAdditionalFloors(ctx context.Context, tx *sql.Tx, listingVersionID int64, floors []listingmodel.WarehouseAdditionalFloorInterface) (err error)

	DeleteListingExchangePlaces(ctx context.Context, tx *sql.Tx, listingVersionID int64) (err error)
	DeleteListingFeatures(ctx context.Context, tx *sql.Tx, listingVersionID int64) (err error)
	DeleteListingGuarantees(ctx context.Context, tx *sql.Tx, listingVersionID int64) (err error)
	DeleteListingFinancingBlockers(ctx context.Context, tx *sql.Tx, listingVersionID int64) (err error)
	DeleteListingWarehouseAdditionalFloors(ctx context.Context, tx *sql.Tx, listingVersionID int64) (err error)

	// Utilities
	GetListingCode(ctx context.Context, tx *sql.Tx) (code uint32, err error)
	GetBaseFeatures(ctx context.Context, tx *sql.Tx) (features []listingmodel.BaseFeatureInterface, err error)
	GetBaseFeaturesByIDs(ctx context.Context, tx *sql.Tx, ids []int64) (map[int64]listingmodel.BaseFeatureInterface, error)

	// GetListingForEndUpdate retrieves comprehensive listing version data for validation flows.
	//
	// This method fetches a specific listing version and returns aggregated data including version metadata,
	// property details, and satellite entity counts. Used primarily by end-update and promote flows.
	//
	// Parameters:
	//   - ctx: Context with request/trace information
	//   - tx: Active database transaction
	//   - versionID: The ID of the listing_version record (NOT the listing_identity_id)
	//
	// Returns:
	//   - ListingEndUpdateData: Aggregated data struct with ListingID populated as listing_identity_id
	//   - error: sql.ErrNoRows if version not found, or infrastructure error
	GetListingForEndUpdate(ctx context.Context, tx *sql.Tx, versionID int64) (ListingEndUpdateData, error)

	ListListings(ctx context.Context, tx *sql.Tx, filter ListListingsFilter) (ListListingsResult, error)

	// New methods for version workflow
	CheckActiveListingExists(ctx context.Context, tx *sql.Tx, userID int64) (bool, error)
	CheckDuplicity(ctx context.Context, tx *sql.Tx, criteria listingmodel.DuplicityCriteria) (bool, error)
	GetListingVersionByAddress(ctx context.Context, tx *sql.Tx, zipCode, number string) (listingmodel.ListingInterface, error)
	GetActiveListingVersion(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (listingmodel.ListingInterface, error)
	GetPreviousActiveVersionStatus(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (listingmodel.ListingStatus, error)
	UpdateListingVersion(ctx context.Context, tx *sql.Tx, version listingmodel.ListingVersionInterface) error
	CloneListingVersionSatellites(ctx context.Context, tx *sql.Tx, sourceVersionID, targetVersionID int64) error

	ListCatalogValues(ctx context.Context, tx *sql.Tx, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error)
	GetCatalogValueByID(ctx context.Context, tx *sql.Tx, category string, id uint8) (listingmodel.CatalogValueInterface, error)
	GetCatalogValueBySlug(ctx context.Context, tx *sql.Tx, category, slug string) (listingmodel.CatalogValueInterface, error)
	GetCatalogValueByNumeric(ctx context.Context, tx *sql.Tx, category string, numericValue uint8) (listingmodel.CatalogValueInterface, error)
	GetNextCatalogValueID(ctx context.Context, tx *sql.Tx, category string) (uint8, error)
	GetNextCatalogNumericValue(ctx context.Context, tx *sql.Tx, category string) (uint8, error)
	CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	UpdateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	SoftDeleteCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8) error

	UpdateListingStatus(ctx context.Context, tx *sql.Tx, listingID int64, newStatus listingmodel.ListingStatus, expectedCurrent listingmodel.ListingStatus) error
	UpdateOwnerResponseStats(ctx context.Context, tx *sql.Tx, identityID int64, deltaSeconds int64, respondedAt time.Time) error
}

type ListingIdentityRecord struct {
	ID              int64
	UUID            string
	UserID          int64
	Code            uint32
	ActiveVersionID sql.NullInt64
	Deleted         bool
}

type ListListingVersionsFilter struct {
	ListingIdentityID int64
	IncludeDeleted    bool
}

type ListingVersionSummary struct {
	Version  listingmodel.ListingVersionInterface
	IsActive bool
}

type ListListingsFilter struct {
	// Pagination
	Page  int // 1-indexed page number (default: 1)
	Limit int // Items per page (default: 20, max: 100)

	// Sorting
	SortBy    string // Field to sort by: id, status (default: id)
	SortOrder string // Sort direction: asc, desc (default: desc)

	// Filters
	Status             *listingmodel.ListingStatus // Optional status filter
	Code               *uint32                     // Optional exact code filter
	Title              string                      // Optional wildcard title/description search
	ZipCode            string                      // Optional wildcard zip code filter
	City               string                      // Optional wildcard city filter
	Neighborhood       string                      // Optional wildcard neighborhood filter
	UserID             *int64                      // Optional owner user ID filter
	MinSellPrice       *float64                    // Optional minimum sell price
	MaxSellPrice       *float64                    // Optional maximum sell price
	MinRentPrice       *float64                    // Optional minimum rent price
	MaxRentPrice       *float64                    // Optional maximum rent price
	MinLandSize        *float64                    // Optional minimum land size (sq meters)
	MaxLandSize        *float64                    // Optional maximum land size (sq meters)
	IncludeAllVersions bool                        // false: active only (default); true: all versions
}

type ListListingsResult struct {
	Records []ListingRecord
	Total   int64
}

type ListingRecord struct {
	Listing listingmodel.ListingInterface
}

// ListingEndUpdateData aggregates the raw values needed to validate the end-update flow.
//
// IMPORTANT: ListingID field contains the listing_identity_id (from listing_versions.listing_identity_id),
// NOT the listing_versions.id. This enables version-to-identity validation in service layer.
type ListingEndUpdateData struct {
	ListingID                  int64 // listing_identity_id from listing_versions table
	UserID                     int64
	Status                     listingmodel.ListingStatus
	Code                       uint32
	Version                    uint8
	ZipCode                    string
	Street                     sql.NullString
	Number                     sql.NullString
	Complex                    sql.NullString
	City                       sql.NullString
	State                      sql.NullString
	Title                      sql.NullString
	ListingType                globalmodel.PropertyType
	Owner                      sql.NullInt16
	Buildable                  sql.NullFloat64
	Delivered                  sql.NullInt16
	WhoLives                   sql.NullInt16
	Description                sql.NullString
	Transaction                sql.NullInt16
	Visit                      sql.NullInt16
	Accompanying               sql.NullInt16
	AnnualTax                  sql.NullFloat64
	MonthlyTax                 sql.NullFloat64
	AnnualGroundRent           sql.NullFloat64
	MonthlyGroundRent          sql.NullFloat64
	Exchange                   sql.NullInt16
	ExchangePercentual         sql.NullFloat64
	SaleNet                    sql.NullFloat64
	RentNet                    sql.NullFloat64
	Condominium                sql.NullFloat64
	LandSize                   sql.NullFloat64
	Corner                     sql.NullInt16
	TenantName                 sql.NullString
	TenantPhone                sql.NullString
	TenantEmail                sql.NullString
	Financing                  sql.NullInt16
	CompletionForecast         sql.NullString
	LandBlock                  sql.NullString
	LandLot                    sql.NullString
	LandFront                  sql.NullFloat64
	LandSide                   sql.NullFloat64
	LandBack                   sql.NullFloat64
	LandTerrainType            sql.NullInt16
	HasKmz                     sql.NullInt16
	KmzFile                    sql.NullString
	BuildingFloors             sql.NullInt16
	UnitTower                  sql.NullString
	UnitFloor                  sql.NullString
	UnitNumber                 sql.NullString
	WarehouseManufacturingArea sql.NullFloat64
	WarehouseSector            sql.NullInt16
	WarehouseHasPrimaryCabin   sql.NullInt16
	WarehouseCabinKva          sql.NullString
	WarehouseGroundFloor       sql.NullInt16
	WarehouseFloorResistance   sql.NullFloat64
	WarehouseZoning            sql.NullString
	WarehouseHasOfficeArea     sql.NullInt16
	WarehouseOfficeArea        sql.NullFloat64
	StoreHasMezzanine          sql.NullInt16
	StoreMezzanineArea         sql.NullFloat64
	FeaturesCount              int
	ExchangePlacesCount        int
	FinancingBlockersCount     int
	GuaranteesCount            int
}
