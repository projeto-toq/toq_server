package listingservices

import (
	"context"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

type listingService struct {
	listingRepository listingrepository.ListingRepoPortInterface
	csi               complexservices.ComplexServiceInterface
	gsi               globalservice.GlobalServiceInterface
	gcs               storageport.CloudStoragePortInterface
}

func NewListingService(
	lr listingrepository.ListingRepoPortInterface,
	csi complexservices.ComplexServiceInterface,
	gsi globalservice.GlobalServiceInterface,
	gcs storageport.CloudStoragePortInterface,

) ListingServiceInterface {
	return &listingService{
		listingRepository: lr,
		csi:               csi,
		gsi:               gsi,
		gcs:               gcs,
	}
}

type ListingServiceInterface interface {
	GetOptions(ctx context.Context, zipCode string, number string) (types []listingmodel.PropertyTypeOption, err error)
	GetBaseFeatures(ctx context.Context) (features []listingmodel.BaseFeatureInterface, err error)
	StartListing(ctx context.Context, zipCode string, number string, propertyType globalmodel.PropertyType) (listing listingmodel.ListingInterface, err error)
	UpdateListing(ctx context.Context, listing listingmodel.ListingInterface) (err error)
	GetAllListingsByUser(ctx context.Context, userID int64) (listings []listingmodel.ListingInterface, err error)
	GetAllOffersByUser(ctx context.Context, userID int64) (offers []listingmodel.OfferInterface, err error)
	GetAllVisitsByUser(ctx context.Context, userID int64) (listings []listingmodel.VisitInterface, err error)
	ApproveOffer(ctx context.Context, offerID int64) (err error)
	RejectOffer(ctx context.Context, offerID int64) (err error)
	ApproveVisit(ctx context.Context, visitID int64) (err error)
	RejectVisit(ctx context.Context, visitID int64) (err error)
	DeleteListing(ctx context.Context, listingID int64) (err error)
	CancelOffer(ctx context.Context, offerID int64) (err error)
	CancelVisit(ctx context.Context, visitID int64) (err error)
}
