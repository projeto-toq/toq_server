package grpcuserport

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (uh *UserHandler) ConfirmEmailChange(ctx context.Context, in *pb.ConfirmEmailChangeRequest) (response *pb.ConfirmEmailChangeResponse, err error) {
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

	tokens, err := uh.service.ConfirmEmailChange(ctx, infos.ID, in.GetCode())
	if err != nil {
		return
	}

	slog.Info("Email change confirmed", "userID", infos.ID, "accessToken", tokens.AccessToken)

	response = &pb.ConfirmEmailChangeResponse{
		Tokens: &pb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}

	return
}
