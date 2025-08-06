package usermodel

type InviteInterface interface {
	GetID() int64
	SetID(int64)
	GetAgencyID() int64
	SetAgencyID(int64)
	GetPhoneNumber() string
	SetPhoneNumber(string)
}

func NewInvite() InviteInterface {
	return &invite{}
}
