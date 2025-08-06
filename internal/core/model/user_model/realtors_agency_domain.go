package usermodel

type realtorsAgency struct {
	id        int64
	agencyID  int64
	realtorID int64
}

func (r *realtorsAgency) GetID() int64 {
	return r.id
}

func (r *realtorsAgency) SetID(id int64) {
	r.id = id
}

func (r *realtorsAgency) GetAgencyID() int64 {
	return r.agencyID
}

func (r *realtorsAgency) SetAgencyID(agencyID int64) {
	r.agencyID = agencyID
}

func (r *realtorsAgency) GetRealtorID() int64 {
	return r.realtorID
}

func (r *realtorsAgency) SetRealtorID(realtorID int64) {
	r.realtorID = realtorID
}
