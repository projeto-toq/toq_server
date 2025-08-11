package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (response *pb.UpdateProfileResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	user, err := validateAndSetInfoForUpdateProfile(in.GetUser())
	if err != nil {
		return
	}

	user.SetID(infos.ID)

	err = uh.service.UpdateProfile(ctx, user)
	if err != nil {
		return
	}

	return &pb.UpdateProfileResponse{}, nil
}

func validateAndSetInfoForUpdateProfile(in *pb.UpdateUser) (user usermodel.UserInterface, err error) { // CEP validation enhancement planned
	user = usermodel.NewUser()
	user.SetZipCode(converters.RemoveAllButDigits(in.GetZipCode()))
	user.SetStreet(converters.TrimSpaces(in.GetStreet()))
	user.SetNumber(converters.TrimSpaces(in.GetNumber()))
	user.SetComplement(converters.TrimSpaces(in.GetComplement()))
	user.SetNeighborhood(converters.TrimSpaces(in.GetNeighborhood()))
	user.SetCity(converters.TrimSpaces(in.GetCity()))
	user.SetState(converters.TrimSpaces(in.GetState()))
	user.SetNickName(converters.TrimSpaces(in.GetNickName()))
	born_at, err := converters.StrngToTime(in.GetBornAT())
	if err != nil {
		return
	}
	user.SetBornAt(born_at)

	return
}
