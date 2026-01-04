package entity

// ServiceArea models photographer_service_areas.
// Columns: id (PK, NOT NULL), photographer_user_id (NOT NULL, FK users.id), city (VARCHAR NOT NULL), state (VARCHAR NOT NULL).
type ServiceArea struct {
	ID                 uint64 // photographer_service_areas.id
	PhotographerUserID uint64 // photographer_service_areas.photographer_user_id (NOT NULL)
	City               string // photographer_service_areas.city (NOT NULL)
	State              string // photographer_service_areas.state (NOT NULL)
}
