package cepmodel

type CEPInterface interface {
	GetCep() string
	SetCep(cep string)
	GetStreet() string
	SetStreet(street string)
	GetComplement() string
	SetComplement(complement string)
	GetNeighborhood() string
	SetNeighborhood(neighborhood string)
	GetCity() string
	SetCity(city string)
	GetState() string
	SetState(state string)
}

func NewCEP() CEPInterface {
	return &cep{}
}
