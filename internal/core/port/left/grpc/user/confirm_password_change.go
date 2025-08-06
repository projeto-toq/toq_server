package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (uh *UserHandler) ConfirmPasswordChange(ctx context.Context, in *pb.ConfirmPasswordChangeRequest) (response *pb.ConfirmPasswordChangeResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	nationalID := converters.RemoveAllButDigits(in.GetNationalID())
	err = validators.ValidateCode(in.GetCode())
	if err != nil {
		return
	}
	password := converters.TrimSpaces(in.GetNewPassword())

	err = uh.service.ConfirmPasswordChange(ctx, nationalID, password, in.GetCode())
	if err != nil {
		return
	}

	return
}
