package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) RequestEmailChange(ctx context.Context, in *pb.RequestEmailChangeRequest) (response *pb.RequestEmailChangeResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	eMail := converters.TrimSpaces(in.GetNewEmail())
	err = validators.ValidateEmail(eMail)
	if err != nil {
		return
	}

	err = uh.service.RequestEmailChange(ctx, infos.ID, eMail)
	if err != nil {
		return
	}

	return
}
