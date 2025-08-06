package usermodel

type RealtorsAgencyInterface interface {
	GetID() int64
	SetID(int64)
	GetAgencyID() int64
	SetAgencyID(int64)
	GetRealtorID() int64
	SetRealtorID(int64)
}

func NewRealtorsAgency() RealtorsAgencyInterface {
	return &realtorsAgency{}
}
