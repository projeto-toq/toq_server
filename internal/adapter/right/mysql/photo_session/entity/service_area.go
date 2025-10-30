package entity

// ServiceArea represents a row from photographer_service_areas.
type ServiceArea struct {
	ID                 uint64
	PhotographerUserID uint64
	City               string
	State              string
}
