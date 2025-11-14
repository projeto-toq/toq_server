package listingmodel

// WarehouseAdditionalFloorInterface represents an additional floor for warehouse properties.
// Each floor has a name, order (for sorting/display), and height in meters.
type WarehouseAdditionalFloorInterface interface {
	ID() int64
	SetID(id int64)
	ListingVersionID() int64
	SetListingVersionID(listingVersionID int64)
	FloorName() string
	SetFloorName(floorName string)
	FloorOrder() int
	SetFloorOrder(floorOrder int)
	FloorHeight() float64
	SetFloorHeight(floorHeight float64)
}
