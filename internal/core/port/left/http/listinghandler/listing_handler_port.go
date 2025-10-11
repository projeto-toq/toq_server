package listinghandler

import "github.com/gin-gonic/gin"

// ListingHandlerPort define a interface para handlers de listing
type ListingHandlerPort interface {
	// Listing management
	GetAllListings(c *gin.Context)
	StartListing(c *gin.Context)
	SearchListing(c *gin.Context)
	PostOptions(c *gin.Context)
	GetBaseFeatures(c *gin.Context)
	GetFavoriteListings(c *gin.Context)
	GetListing(c *gin.Context)
	UpdateListing(c *gin.Context)
	DeleteListing(c *gin.Context)
	EndUpdateListing(c *gin.Context)
	GetListingStatus(c *gin.Context)
	ApproveListing(c *gin.Context)
	RejectListing(c *gin.Context)
	SuspendListing(c *gin.Context)
	ReleaseListing(c *gin.Context)
	CopyListing(c *gin.Context)
	ShareListing(c *gin.Context)
	AddFavoriteListing(c *gin.Context)
	RemoveFavoriteListing(c *gin.Context)

	// Visit management
	RequestVisit(c *gin.Context)
	GetVisits(c *gin.Context)
	GetAllVisits(c *gin.Context)
	CancelVisit(c *gin.Context)
	ConfirmVisitDone(c *gin.Context)
	ApproveVisting(c *gin.Context)
	RejectVisting(c *gin.Context)

	// Offer management
	CreateOffer(c *gin.Context)
	GetOffers(c *gin.Context)
	GetAllOffers(c *gin.Context)
	UpdateOffer(c *gin.Context)
	CancelOffer(c *gin.Context)
	SendOffer(c *gin.Context)
	ApproveOffer(c *gin.Context)
	RejectOffer(c *gin.Context)

	// Evaluation
	EvaluateRealtor(c *gin.Context)
	EvaluateOwner(c *gin.Context)
}
