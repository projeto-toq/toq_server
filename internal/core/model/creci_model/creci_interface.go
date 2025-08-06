package crecimodel

import "time"

type CreciInterface interface {
	GetCreciNumber() string
	SetCreciNumber(string)
	GetCreciState() string
	SetCreciState(string)
	GetCreciValidity() time.Time
	SetCreciValidity(time.Time)
}

func NewCreci() CreciInterface {
	return &creci{}
}
