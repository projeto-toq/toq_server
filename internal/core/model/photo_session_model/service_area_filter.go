package photosessionmodel

// ServiceAreaFilter represents filtering options for service area listing.
type ServiceAreaFilter struct {
	City   *string
	State  *string
	Offset int
	Limit  int
}
