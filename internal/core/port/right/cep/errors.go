package cepport

import "errors"

var (
	ErrInvalid     = errors.New("cep invalid")
	ErrNotFound    = errors.New("cep not found")
	ErrRateLimited = errors.New("cep rate limited")
	ErrInfra       = errors.New("cep infra error")
)
