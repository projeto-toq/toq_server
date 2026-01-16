package listinghandler

import "github.com/gin-gonic/gin"

// ListingHandlerPort define a interface para handlers de listing
type ListingHandlerPort interface {
	// Listing management
	ListListings(c *gin.Context)
	ListComplexes(c *gin.Context)
	StartListing(c *gin.Context)
	CreateDraftVersion(c *gin.Context)
	PostOptions(c *gin.Context)
	GetBaseFeatures(c *gin.Context)
	ListCatalogValues(c *gin.Context)
	GetFavoriteListings(c *gin.Context)
	AddFavoriteListing(c *gin.Context)
	RemoveFavoriteListing(c *gin.Context)
	ListPhotographerSlots(c *gin.Context)
	ReservePhotoSession(c *gin.Context)
	CancelPhotoSession(c *gin.Context)
	GetListing(c *gin.Context)
	UpdateListing(c *gin.Context)
	PromoteListingVersion(c *gin.Context)
	DiscardDraftVersion(c *gin.Context)
	ListListingVersions(c *gin.Context)
	ChangeListingStatus(c *gin.Context)
}
