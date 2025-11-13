package listingrepository

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

type ListingRepoPortInterface interface {
	CreateListing(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) (err error)
	CreateExchangePlace(ctx context.Context, tx *sql.Tx, place listingmodel.ExchangePlaceInterface) (err error)
	CreateFeature(ctx context.Context, tx *sql.Tx, feature listingmodel.FeatureInterface) (err error)
	CreateGuarantee(ctx context.Context, tx *sql.Tx, guarantee listingmodel.GuaranteeInterface) (err error)
	CreateFinancingBlocker(ctx context.Context, tx *sql.Tx, blocker listingmodel.FinancingBlockerInterface) (err error)

	UpdateListing(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) (err error)
	UpdateExchangePlaces(ctx context.Context, tx *sql.Tx, listingID int64, places []listingmodel.ExchangePlaceInterface) (err error)
	UpdateFeatures(ctx context.Context, tx *sql.Tx, listingID int64, features []listingmodel.FeatureInterface) (err error)
	UpdateGuarantees(ctx context.Context, tx *sql.Tx, listingID int64, guarantees []listingmodel.GuaranteeInterface) (err error)
	UpdateFinancingBlockers(ctx context.Context, tx *sql.Tx, listingID int64, blockers []listingmodel.FinancingBlockerInterface) (err error)

	DeleteListingExchangePlaces(ctx context.Context, tx *sql.Tx, listingID int64) (err error)
	DeleteListingFeatures(ctx context.Context, tx *sql.Tx, listingID int64) (err error)
	DeleteListingGuarantees(ctx context.Context, tx *sql.Tx, listingID int64) (err error)
	DeleteListingFinancingBlockers(ctx context.Context, tx *sql.Tx, listingID int64) (err error)

	GetListingCode(ctx context.Context, tx *sql.Tx) (code uint32, err error)
	GetBaseFeatures(ctx context.Context, tx *sql.Tx) (features []listingmodel.BaseFeatureInterface, err error)
	GetBaseFeaturesByIDs(ctx context.Context, tx *sql.Tx, ids []int64) (map[int64]listingmodel.BaseFeatureInterface, error)
	GetListingByZipNumber(ctx context.Context, tx *sql.Tx, zip string, number string) (listing listingmodel.ListingInterface, err error)
	GetListingByID(ctx context.Context, tx *sql.Tx, listingID int64) (listing listingmodel.ListingInterface, err error)
	GetListingForEndUpdate(ctx context.Context, tx *sql.Tx, listingID int64) (ListingEndUpdateData, error)
	ListListings(ctx context.Context, tx *sql.Tx, filter ListListingsFilter) (ListListingsResult, error)

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
}

type ListListingsFilter struct {
	Page         int
	Limit        int
	Status       *listingmodel.ListingStatus
	Code         *uint32
	Title        string
	ZipCode      string
	City         string
	Neighborhood string
	UserID       *int64
	MinSellPrice *float64
	MaxSellPrice *float64
	MinRentPrice *float64
	MaxRentPrice *float64
	MinLandSize  *float64
	MaxLandSize  *float64
}

type ListListingsResult struct {
	Records []ListingRecord
	Total   int64
}

type ListingRecord struct {
	Listing listingmodel.ListingInterface
}

// ListingEndUpdateData aggregates the raw values needed to validate the end-update flow.
type ListingEndUpdateData struct {
	ListingID              int64
	UserID                 int64
	Status                 listingmodel.ListingStatus
	Code                   uint32
	Version                uint8
	ZipCode                string
	Street                 sql.NullString
	Number                 sql.NullString
	City                   sql.NullString
	State                  sql.NullString
	Title                  sql.NullString
	ListingType            globalmodel.PropertyType
	Owner                  sql.NullInt16
	Buildable              sql.NullFloat64
	Delivered              sql.NullInt16
	WhoLives               sql.NullInt16
	Description            sql.NullString
	Transaction            sql.NullInt16
	Visit                  sql.NullInt16
	Accompanying           sql.NullInt16
	AnnualTax              sql.NullFloat64
	MonthlyTax             sql.NullFloat64
	AnnualGroundRent       sql.NullFloat64
	MonthlyGroundRent      sql.NullFloat64
	Exchange               sql.NullInt16
	ExchangePercentual     sql.NullFloat64
	SaleNet                sql.NullFloat64
	RentNet                sql.NullFloat64
	Condominium            sql.NullFloat64
	LandSize               sql.NullFloat64
	Corner                 sql.NullInt16
	TenantName             sql.NullString
	TenantPhone            sql.NullString
	TenantEmail            sql.NullString
	Financing              sql.NullInt16
	FeaturesCount          int
	ExchangePlacesCount    int
	FinancingBlockersCount int
	GuaranteesCount        int
}
