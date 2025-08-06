package creciport

import (
	"context"

	crecimodel "github.com/giulio-alfieri/toq_server/internal/core/model/creci_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

type CreciPortInterface interface {
	Open(context.Context) error
	Close()
	ValidateCreciNumber(ctx context.Context, realtor usermodel.UserInterface) (creci crecimodel.CreciInterface, err error)
	ValidateFaceMatch(ctx context.Context, realtor usermodel.UserInterface) (match bool, err error)
}
