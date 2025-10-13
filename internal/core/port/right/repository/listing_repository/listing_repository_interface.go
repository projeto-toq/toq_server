package listingrepository

import (
	"context"
	"database/sql"

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
	GetListingByZipNumber(ctx context.Context, tx *sql.Tx, zip string, number string) (listing listingmodel.ListingInterface, err error)
	GetListingByID(ctx context.Context, tx *sql.Tx, listingID int64) (listing listingmodel.ListingInterface, err error)

	ListCatalogValues(ctx context.Context, tx *sql.Tx, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error)
	GetCatalogValueByID(ctx context.Context, tx *sql.Tx, category string, id uint8) (listingmodel.CatalogValueInterface, error)
	GetCatalogValueBySlug(ctx context.Context, tx *sql.Tx, category, slug string) (listingmodel.CatalogValueInterface, error)
	GetNextCatalogValueID(ctx context.Context, tx *sql.Tx, category string) (uint8, error)
	CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	UpdateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	SoftDeleteCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8) error
}
