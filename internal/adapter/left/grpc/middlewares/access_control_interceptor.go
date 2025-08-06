package middlewares

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AccessControlInterceptor(ctx context.Context, memCache *cache.CacheInterface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		cache := *memCache
		infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

		//lógica de controle de acesso
		allowed, valid, err := cache.Get(ctx, info.FullMethod, infos.Role)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal server error")
		}
		if (!valid || !allowed) && infos.Role != usermodel.RoleRoot {
			return nil, status.Error(codes.PermissionDenied, "user not authorized for this action")
		}

		// Handle the request
		resp, err := handler(ctx, req)

		return resp, err
	}
}
