package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetRealtorsByAgency(ctx context.Context, in *pb.GetRealtorsByAgencyRequest) (out *pb.GetRealtorsByAgencyResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	realtors, err := uh.service.GetRealtorsByAgency(ctx, infos.ID)
	if err != nil {
		return
	}
	out = &pb.GetRealtorsByAgencyResponse{
		Realtors: realtorToPBUser(realtors),
	}

	return
}

func realtorToPBUser(realtors []usermodel.UserInterface) (users []*pb.User) {
	for _, realtor := range realtors {
		user := pb.User{
			Id:            realtor.GetID(),
			FullName:      realtor.GetFullName(),
			NickName:      realtor.GetNickName(),
			NationalID:    realtor.GetNationalID(),
			CreciNumber:   realtor.GetCreciNumber(),
			CreciState:    realtor.GetCreciState(),
			CreciValidity: realtor.GetCreciValidity().String(),
			BornAT:        realtor.GetBornAt().String(),
			PhoneNumber:   realtor.GetPhoneNumber(),
			Email:         realtor.GetEmail(),
			ZipCode:       realtor.GetZipCode(),
			Street:        realtor.GetStreet(),
			Number:        realtor.GetNumber(),
			Complement:    realtor.GetComplement(),
			Neighborhood:  realtor.GetNeighborhood(),
			City:          realtor.GetCity(),
			State:         realtor.GetState(),
		}
		users = append(users, &user)
	}
	return
}
