package grpcuserport

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	userHandlerconverters "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/converters"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (response *pb.GetProfileResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	user, err := uh.service.GetProfile(ctx, infos.ID)
	if err != nil {
		return nil, err
	}

	// Generate signed URL for photo download
	photoURL, err := uh.service.GeneratePhotoDownloadURL(ctx, infos.ID)
	if err != nil {
		slog.Error("failed to generate signed URL for photo", "error", err, "userID", infos.ID)
		photoURL = "" // Continue with empty photo URL
	}

	return userHandlerconverters.UserDomainToProfileResponse(user, photoURL)
}
