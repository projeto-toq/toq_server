package cnpjport

import "errors"

var (
	ErrInvalid     = errors.New("cnpj invalid")
	ErrNotFound    = errors.New("cnpj not found")
	ErrRateLimited = errors.New("cnpj rate limited")
	ErrInfra       = errors.New("cnpj infra error")
)
