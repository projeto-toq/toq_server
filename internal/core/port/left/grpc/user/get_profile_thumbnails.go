package grpcuserport

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetProfileThumbnails(ctx context.Context, req *pb.GetProfileThumbnailsRequest) (*pb.GetProfileThumbnailsResponse, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Extrair informações do usuário do token
	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	// Buscar todas as URLs dos thumbnails
	thumbnails, err := uh.service.GetProfileThumbnails(ctx, infos.ID)
	if err != nil {
		slog.Error("failed to get profile thumbnails", "error", err, "userID", infos.ID)
		return nil, err
	}

	return &pb.GetProfileThumbnailsResponse{
		OriginalUrl: thumbnails.OriginalURL,
		SmallUrl:    thumbnails.SmallURL,
		MediumUrl:   thumbnails.MediumURL,
		LargeUrl:    thumbnails.LargeURL,
	}, nil
}
