package grpcuserport

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) ConfirmPhoneChange(ctx context.Context, in *pb.ConfirmPhoneChangeRequest) (response *pb.ConfirmPhoneChangeResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	err = validators.ValidateCode(in.GetCode())
	if err != nil {
		return
	}

	tokens, err := uh.service.ConfirmPhoneChange(ctx, infos.ID, in.GetCode())
	if err != nil {
		return
	}

	slog.Info("Phone change confirmed", "userID", infos.ID, "accessToken", tokens.AccessToken)

	response = &pb.ConfirmPhoneChangeResponse{
		Tokens: &pb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}

	return
}
