package photosessionmodel

import "time"

// PhotographerTimeOffInterface defines contract for photographer time-off periods.
type PhotographerTimeOffInterface interface {
	ID() uint64
	SetID(value uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(value uint64)
	StartDate() time.Time
	SetStartDate(value time.Time)
	EndDate() time.Time
	SetEndDate(value time.Time)
	Reason() *string
	SetReason(value *string)
}

// NewPhotographerTimeOff instantiates an empty time-off entity.
func NewPhotographerTimeOff() PhotographerTimeOffInterface {
	return &photographerTimeOff{}
}
