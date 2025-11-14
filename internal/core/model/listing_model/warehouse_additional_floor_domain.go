package listingmodel

type warehouseAdditionalFloor struct {
	id               int64
	listingVersionID int64
	floorName        string
	floorOrder       int
	floorHeight      float64
}

// NewWarehouseAdditionalFloor creates a new warehouse additional floor instance.
func NewWarehouseAdditionalFloor() WarehouseAdditionalFloorInterface {
	return &warehouseAdditionalFloor{}
}

func (w *warehouseAdditionalFloor) ID() int64 {
	return w.id
}

func (w *warehouseAdditionalFloor) SetID(id int64) {
	w.id = id
}

func (w *warehouseAdditionalFloor) ListingVersionID() int64 {
	return w.listingVersionID
}

func (w *warehouseAdditionalFloor) SetListingVersionID(listingVersionID int64) {
	w.listingVersionID = listingVersionID
}

func (w *warehouseAdditionalFloor) FloorName() string {
	return w.floorName
}

func (w *warehouseAdditionalFloor) SetFloorName(floorName string) {
	w.floorName = floorName
}

func (w *warehouseAdditionalFloor) FloorOrder() int {
	return w.floorOrder
}

func (w *warehouseAdditionalFloor) SetFloorOrder(floorOrder int) {
	w.floorOrder = floorOrder
}

func (w *warehouseAdditionalFloor) FloorHeight() float64 {
	return w.floorHeight
}

func (w *warehouseAdditionalFloor) SetFloorHeight(floorHeight float64) {
	w.floorHeight = floorHeight
}
