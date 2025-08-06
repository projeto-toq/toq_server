package usermodel

type invite struct {
	id          int64
	agencyID    int64
	phoneNumber string
}

func (i *invite) GetID() int64 {
	return i.id
}

func (i *invite) SetID(id int64) {
	i.id = id
}

func (i *invite) GetAgencyID() int64 {
	return i.agencyID
}

func (i *invite) SetAgencyID(agencyID int64) {
	i.agencyID = agencyID
}

func (i *invite) GetPhoneNumber() string {
	return i.phoneNumber
}

func (i *invite) SetPhoneNumber(phoneNumber string) {
	i.phoneNumber = phoneNumber
}
