package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
)

// GetListingById handles getting a specific listing by ID
// Service method not implemented yet
func (lh *ListingHandler) GetListingById(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingById service method not implemented yet")
}

// GetListingByComplexId handles getting listings by complex ID
// Service method not implemented yet
func (lh *ListingHandler) GetListingByComplexId(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingByComplexId service method not implemented yet")
}

// ChangeListingStatus handles changing a listing's status
// Service method not implemented yet
func (lh *ListingHandler) ChangeListingStatus(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "ChangeListingStatus service method not implemented yet")
}

// GetUserListingsByStatus handles getting user listings filtered by status
// Service method not implemented yet
func (lh *ListingHandler) GetUserListingsByStatus(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetUserListingsByStatus service method not implemented yet")
}

// AddListingPhotos handles adding photos to a listing
// Service method not implemented yet
func (lh *ListingHandler) AddListingPhotos(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "AddListingPhotos service method not implemented yet")
}

// RemoveListingPhoto handles removing a photo from a listing
// Service method not implemented yet
func (lh *ListingHandler) RemoveListingPhoto(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "RemoveListingPhoto service method not implemented yet")
}

// UpdateListingPhoto handles updating a listing photo
// Service method not implemented yet
func (lh *ListingHandler) UpdateListingPhoto(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "UpdateListingPhoto service method not implemented yet")
}

// GetListingPhotos handles getting photos for a listing
// Service method not implemented yet
func (lh *ListingHandler) GetListingPhotos(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingPhotos service method not implemented yet")
}

// GetListingFeatures handles getting features for a listing
// Service method not implemented yet
func (lh *ListingHandler) GetListingFeatures(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingFeatures service method not implemented yet")
}

// AddListingFeatures handles adding features to a listing
// Service method not implemented yet
func (lh *ListingHandler) AddListingFeatures(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "AddListingFeatures service method not implemented yet")
}

// RemoveListingFeature handles removing a feature from a listing
// Service method not implemented yet
func (lh *ListingHandler) RemoveListingFeature(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "RemoveListingFeature service method not implemented yet")
}

// GetListingComments handles getting comments for a listing
// Service method not implemented yet
func (lh *ListingHandler) GetListingComments(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingComments service method not implemented yet")
}

// AddListingComment handles adding a comment to a listing
// Service method not implemented yet
func (lh *ListingHandler) AddListingComment(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "AddListingComment service method not implemented yet")
}

// GetListingViews handles getting views for a listing
// Service method not implemented yet
func (lh *ListingHandler) GetListingViews(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingViews service method not implemented yet")
}

// AddListingView handles adding a view to a listing
// Service method not implemented yet
func (lh *ListingHandler) AddListingView(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "AddListingView service method not implemented yet")
}

// GetListingFavorites handles getting favorites for a listing
// Service method not implemented yet
func (lh *ListingHandler) GetListingFavorites(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingFavorites service method not implemented yet")
}

// AddListingFavorite handles adding a favorite to a listing
// Service method not implemented yet
func (lh *ListingHandler) AddListingFavorite(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "AddListingFavorite service method not implemented yet")
}

// SearchListing handles searching for listings
// Service method not implemented yet
func (lh *ListingHandler) SearchListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "SearchListing service method not implemented yet")
}

// GetFavoriteListings handles getting favorite listings
// Service method not implemented yet
func (lh *ListingHandler) GetFavoriteListings(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetFavoriteListings service method not implemented yet")
}

// GetListingStatus handles getting listing status
// Service method not implemented yet
func (lh *ListingHandler) GetListingStatus(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetListingStatus service method not implemented yet")
}

// ApproveListing handles approving a listing
// Service method not implemented yet
func (lh *ListingHandler) ApproveListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "ApproveListing service method not implemented yet")
}

// RejectListing handles rejecting a listing
// Service method not implemented yet
func (lh *ListingHandler) RejectListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "RejectListing service method not implemented yet")
}

// SuspendListing handles suspending a listing
// Service method not implemented yet
func (lh *ListingHandler) SuspendListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "SuspendListing service method not implemented yet")
}

// ReleaseListing handles releasing a listing
// Service method not implemented yet
func (lh *ListingHandler) ReleaseListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "ReleaseListing service method not implemented yet")
}

// CopyListing handles copying a listing
// Service method not implemented yet
func (lh *ListingHandler) CopyListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "CopyListing service method not implemented yet")
}

// ShareListing handles sharing a listing
// Service method not implemented yet
func (lh *ListingHandler) ShareListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "ShareListing service method not implemented yet")
}

// AddFavoriteListing handles adding a listing to favorites
// Service method not implemented yet
func (lh *ListingHandler) AddFavoriteListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "AddFavoriteListing service method not implemented yet")
}

// RemoveFavoriteListing handles removing a listing from favorites
// Service method not implemented yet
func (lh *ListingHandler) RemoveFavoriteListing(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "RemoveFavoriteListing service method not implemented yet")
}

// RequestVisit handles requesting a visit
// Service method not implemented yet
func (lh *ListingHandler) RequestVisit(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "RequestVisit service method not implemented yet")
}

// GetVisits handles getting visits for a listing
// Service method not implemented yet
func (lh *ListingHandler) GetVisits(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetVisits service method not implemented yet")
}

// GetAllVisits handles getting all visits
// Service method not implemented yet
func (lh *ListingHandler) GetAllVisits(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetAllVisits service method not implemented yet")
}

// CancelVisit handles canceling a visit
// Service method not implemented yet
func (lh *ListingHandler) CancelVisit(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "CancelVisit service method not implemented yet")
}

// ConfirmVisitDone handles confirming a visit is done
// Service method not implemented yet
func (lh *ListingHandler) ConfirmVisitDone(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "ConfirmVisitDone service method not implemented yet")
}

// ApproveVisting handles approving a visit
// Service method not implemented yet
func (lh *ListingHandler) ApproveVisting(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "ApproveVisting service method not implemented yet")
}

// RejectVisting handles rejecting a visit
// Service method not implemented yet
func (lh *ListingHandler) RejectVisting(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "RejectVisting service method not implemented yet")
}

// CreateOffer handles creating an offer
// Service method not implemented yet
func (lh *ListingHandler) CreateOffer(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "CreateOffer service method not implemented yet")
}

// GetOffers handles getting offers for a listing
// Service method not implemented yet
func (lh *ListingHandler) GetOffers(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetOffers service method not implemented yet")
}

// GetAllOffers handles getting all offers
// Service method not implemented yet
func (lh *ListingHandler) GetAllOffers(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetAllOffers service method not implemented yet")
}

// UpdateOffer handles updating an offer
// Service method not implemented yet
func (lh *ListingHandler) UpdateOffer(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "UpdateOffer service method not implemented yet")
}

// CancelOffer handles canceling an offer
// Service method not implemented yet
func (lh *ListingHandler) CancelOffer(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "CancelOffer service method not implemented yet")
}

// SendOffer handles sending an offer
// Service method not implemented yet
func (lh *ListingHandler) SendOffer(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "SendOffer service method not implemented yet")
}

// ApproveOffer handles approving an offer
// Service method not implemented yet
func (lh *ListingHandler) ApproveOffer(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "ApproveOffer service method not implemented yet")
}

// RejectOffer handles rejecting an offer
// Service method not implemented yet
func (lh *ListingHandler) RejectOffer(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "RejectOffer service method not implemented yet")
}

// EvaluateRealtor handles evaluating a realtor
// Service method not implemented yet
func (lh *ListingHandler) EvaluateRealtor(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "EvaluateRealtor service method not implemented yet")
}

// EvaluateOwner handles evaluating an owner
// Service method not implemented yet
func (lh *ListingHandler) EvaluateOwner(c *gin.Context) {
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "EvaluateOwner service method not implemented yet")
}
