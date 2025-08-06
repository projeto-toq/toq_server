package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"
)

func (uh *UserHandler) RequestPasswordChange(ctx context.Context, in *pb.RequestPasswordChangeRequest) (response *pb.RequestPasswordChangeResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	nationalID := converters.RemoveAllButDigits(in.GetNationalID())

	err = uh.service.RequestPasswordChange(ctx, nationalID)
	if err != nil {
		return
	}

	return
}
