package listingfavoriteconverters

import listingfavoriteentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing_favorite/entities"

// FavoriteEntityToIdentityID extracts the listing identity ID from the favorite entity.
func FavoriteEntityToIdentityID(entity listingfavoriteentity.FavoriteEntity) int64 {
	return entity.ListingIdentityID
}
