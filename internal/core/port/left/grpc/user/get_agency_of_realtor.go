package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/golang/protobuf/ptypes/empty"
)

func (uh *UserHandler) GetAgencyOfRealtor(ctx context.Context, empty *empty.Empty) (out *pb.GetAgencyOfRealtorResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	agency, err := uh.service.GetAgencyOfRealtor(ctx, infos.ID)
	if err != nil {
		return
	}

	out = &pb.GetAgencyOfRealtorResponse{
		Agency: &pb.User{
			Id:           agency.GetID(),
			FullName:     agency.GetFullName(),
			NickName:     agency.GetNickName(),
			NationalID:   agency.GetNationalID(),
			CreciNumber:  agency.GetCreciNumber(),
			CreciState:   agency.GetCreciState(),
			PhoneNumber:  agency.GetPhoneNumber(),
			Email:        agency.GetEmail(),
			ZipCode:      agency.GetZipCode(),
			Street:       agency.GetStreet(),
			Number:       agency.GetNumber(),
			Complement:   agency.GetComplement(),
			Neighborhood: agency.GetNeighborhood(),
			City:         agency.GetCity(),
			State:        agency.GetState(),
		},
	}
	return
}
