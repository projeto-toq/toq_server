package photosessionmodel

// PhotographerServiceAreaInterface defines how service area entities can be manipulated in the domain.
type PhotographerServiceAreaInterface interface {
	ID() uint64
	SetID(id uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(id uint64)
	City() string
	SetCity(city string)
	State() string
	SetState(state string)
}

// NewPhotographerServiceArea constructs an empty PhotographerServiceAreaInterface implementation.
func NewPhotographerServiceArea() PhotographerServiceAreaInterface {
	return &photographerServiceArea{}
}
