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

func (uh *UserHandler) RequestPhoneChange(ctx context.Context, in *pb.RequestPhoneChangeRequest) (response *pb.RequestPhoneChangeResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	newPhone := converters.RemoveAllButDigitsAndPlusSign(in.GetNewPhoneNumber())
	err = validators.ValidateE164(newPhone)
	if err != nil {
		return
	}

	err = uh.service.RequestPhoneChange(ctx, infos.ID, newPhone)
	if err != nil {
		return
	}

	return
}
