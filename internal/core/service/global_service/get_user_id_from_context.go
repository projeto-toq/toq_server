package globalservice

import (
	"context"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUserIDFromContext extrai o ID do usuário do contexto.
// O ID é injetado no contexto pelo interceptor de autenticação.
func (gs *globalService) GetUserIDFromContext(ctx context.Context) (int64, error) {
	userInfos, ok := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)
	if !ok || userInfos.ID == 0 {
		return 0, status.Error(codes.Unauthenticated, "user not found in context")
	}
	return userInfos.ID, nil
}
