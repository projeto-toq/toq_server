package middlewares

import (
	"context"
	"log/slog"

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
			slog.Error("Error getting permissions from cache", "method", info.FullMethod, "role", infos.Role, "error", err)
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		slog.Debug("Permission check result", "method", info.FullMethod, "role", infos.Role, "allowed", allowed, "valid", valid)

		if (!valid || !allowed) && infos.Role != usermodel.RoleRoot {
			slog.Warn("Usuário não tem acesso a este RPC", "method", info.FullMethod, "role", infos.Role, "allowed", allowed, "valid", valid)
			return nil, status.Error(codes.PermissionDenied, "user not authorized for this action")
		}

		// Handle the request
		resp, err := handler(ctx, req)

		return resp, err
	}
}
