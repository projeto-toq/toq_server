package listingfavoriteentity

// FavoriteEntity mirrors the listing_favorites table. Adapter scope only.
type FavoriteEntity struct {
	ID                int64
	ListingIdentityID int64
	UserID            int64
}
