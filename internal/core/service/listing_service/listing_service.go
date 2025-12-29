package listingservices

import (
	"context"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	userrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
)

type listingService struct {
	listingRepository listingrepository.ListingRepoPortInterface
	photoSessionSvc   photosessionservices.PhotoSessionServiceInterface
	userRepository    userrepository.UserRepoPortInterface
	propertyCoverage  propertycoverageservice.PropertyCoverageServiceInterface
	gsi               globalservice.GlobalServiceInterface
	gcs               storageport.CloudStoragePortInterface
	scheduleService   scheduleservices.ScheduleServiceInterface
}

func NewListingService(
	lr listingrepository.ListingRepoPortInterface,
	ps photosessionservices.PhotoSessionServiceInterface,
	ur userrepository.UserRepoPortInterface,
	pcs propertycoverageservice.PropertyCoverageServiceInterface,
	gsi globalservice.GlobalServiceInterface,
	gcs storageport.CloudStoragePortInterface,
	ss scheduleservices.ScheduleServiceInterface,

) ListingServiceInterface {
	return &listingService{
		listingRepository: lr,
		photoSessionSvc:   ps,
		userRepository:    ur,
		propertyCoverage:  pcs,
		gsi:               gsi,
		gcs:               gcs,
		scheduleService:   ss,
	}
}

type ListingServiceInterface interface {
	GetOptions(ctx context.Context, zipCode string, number string) (propertycoveragemodel.ManagedComplexInterface, error)
	GetBaseFeatures(ctx context.Context) (features []listingmodel.BaseFeatureInterface, err error)
	CreateListing(ctx context.Context, input CreateListingInput) (listing listingmodel.ListingInterface, err error)
	CreateDraftVersion(ctx context.Context, input CreateDraftVersionInput) (CreateDraftVersionOutput, error)
	UpdateListing(ctx context.Context, input UpdateListingInput) (err error)
	PromoteListingVersion(ctx context.Context, input PromoteListingVersionInput) error
	ChangeListingStatus(ctx context.Context, input ChangeListingStatusInput) (ChangeListingStatusOutput, error)
	DiscardDraftVersion(ctx context.Context, input DiscardDraftVersionInput) error
	ListListingVersions(ctx context.Context, input ListListingVersionsInput) (ListListingVersionsOutput, error)
	GetAllListingsByUser(ctx context.Context, userID int64) (listings []listingmodel.ListingInterface, err error)
	ListListings(ctx context.Context, input ListListingsInput) (ListListingsOutput, error)
	GetAllOffersByUser(ctx context.Context, userID int64) (offers []listingmodel.OfferInterface, err error)
	GetAllVisitsByUser(ctx context.Context, userID int64) (listings []listingmodel.VisitInterface, err error)
	ApproveOffer(ctx context.Context, offerID int64) (err error)
	RejectOffer(ctx context.Context, offerID int64) (err error)
	ApproveVisit(ctx context.Context, visitID int64) (err error)
	RejectVisit(ctx context.Context, visitID int64) (err error)
	DeleteListing(ctx context.Context, listingID int64) (err error)
	CancelOffer(ctx context.Context, offerID int64) (err error)
	CancelVisit(ctx context.Context, visitID int64) (err error)
	ListCatalogValues(ctx context.Context, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error)
	CreateCatalogValue(ctx context.Context, input CreateCatalogValueInput) (listingmodel.CatalogValueInterface, error)
	UpdateCatalogValue(ctx context.Context, input UpdateCatalogValueInput) (listingmodel.CatalogValueInterface, error)
	DeleteCatalogValue(ctx context.Context, category string, id uint8) error
	RestoreCatalogValue(ctx context.Context, input RestoreCatalogValueInput) (listingmodel.CatalogValueInterface, error)
	GetCatalogValueDetail(ctx context.Context, category string, id uint8) (listingmodel.CatalogValueInterface, error)
	ListPhotographerSlots(ctx context.Context, input ListPhotographerSlotsInput) (ListPhotographerSlotsOutput, error)
	ReservePhotoSession(ctx context.Context, input ReservePhotoSessionInput) (ReservePhotoSessionOutput, error)
	CancelPhotoSession(ctx context.Context, input CancelPhotoSessionInput) error
	GetListingDetail(ctx context.Context, listingIdentityId int64) (ListingDetailOutput, error)
}
