package listingentity

// EntityWarehouseAdditionalFloor represents a warehouse additional floor record in the database.
type EntityWarehouseAdditionalFloor struct {
	ID               int64
	ListingVersionID int64
	FloorName        string
	FloorOrder       int
	FloorHeight      float64
}
