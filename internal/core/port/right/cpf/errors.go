package cpfport

import "errors"

var (
	ErrInvalidInput     = errors.New("cpf invalid input")
	ErrBirthDateInvalid = errors.New("cpf birth date invalid")
	ErrNotFound         = errors.New("cpf not found")
	ErrStatusIrregular  = errors.New("cpf status irregular")
	ErrDataMismatch     = errors.New("cpf data mismatch")
	ErrRateLimited      = errors.New("cpf rate limited")
	ErrInfra            = errors.New("cpf infra error")
)
