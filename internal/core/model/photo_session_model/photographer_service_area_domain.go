package photosessionmodel

type photographerServiceArea struct {
	id                 uint64
	photographerUserID uint64
	city               string
	state              string
}

func (a *photographerServiceArea) ID() uint64 { return a.id }

func (a *photographerServiceArea) SetID(id uint64) { a.id = id }

func (a *photographerServiceArea) PhotographerUserID() uint64 { return a.photographerUserID }

func (a *photographerServiceArea) SetPhotographerUserID(id uint64) { a.photographerUserID = id }

func (a *photographerServiceArea) City() string { return a.city }

func (a *photographerServiceArea) SetCity(city string) { a.city = city }

func (a *photographerServiceArea) State() string { return a.state }

func (a *photographerServiceArea) SetState(state string) { a.state = state }
