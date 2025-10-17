package listingservices

import (
	"context"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	photosessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/photo_session_repository"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

type listingService struct {
	listingRepository listingrepository.ListingRepoPortInterface
	photoSessionRepo  photosessionrepository.PhotoSessionRepositoryInterface
	csi               complexservices.ComplexServiceInterface
	gsi               globalservice.GlobalServiceInterface
	gcs               storageport.CloudStoragePortInterface
}

func NewListingService(
	lr listingrepository.ListingRepoPortInterface,
	ps photosessionrepository.PhotoSessionRepositoryInterface,
	csi complexservices.ComplexServiceInterface,
	gsi globalservice.GlobalServiceInterface,
	gcs storageport.CloudStoragePortInterface,

) ListingServiceInterface {
	return &listingService{
		listingRepository: lr,
		photoSessionRepo:  ps,
		csi:               csi,
		gsi:               gsi,
		gcs:               gcs,
	}
}

type ListingServiceInterface interface {
	GetOptions(ctx context.Context, zipCode string, number string) (types []listingmodel.PropertyTypeOption, err error)
	GetBaseFeatures(ctx context.Context) (features []listingmodel.BaseFeatureInterface, err error)
	StartListing(ctx context.Context, input StartListingInput) (listing listingmodel.ListingInterface, err error)
	UpdateListing(ctx context.Context, input UpdateListingInput) (err error)
	EndUpdateListing(ctx context.Context, input EndUpdateListingInput) error
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
	ListPhotographerSlots(ctx context.Context, input ListPhotographerSlotsInput) (ListPhotographerSlotsOutput, error)
	ReservePhotoSession(ctx context.Context, input ReservePhotoSessionInput) (ReservePhotoSessionOutput, error)
	ConfirmPhotoSession(ctx context.Context, input ConfirmPhotoSessionInput) (ConfirmPhotoSessionOutput, error)
	CancelPhotoSession(ctx context.Context, input CancelPhotoSessionInput) error
	GetListingDetail(ctx context.Context, listingID int64) (ListingDetailOutput, error)
}
