package userservices

import (
	"context"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (us *userService) CleanAgencyPending(ctx context.Context, agency usermodel.UserInterface) (err error) {

	realtors, err := us.GetRealtorsByAgency(ctx, agency.GetID())
	if err != nil {
		return
	}

	for _, realtor := range realtors {
		err = us.DeleteRealtorOfAgency(ctx, agency.GetID(), realtor.GetID())
		if err != nil {
			return
		}
	}

	return
}
