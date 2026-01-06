package globalservice

import (
	"context"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserIDFromContext extracts the authenticated user ID from the context.
// The authentication middleware injects the token payload using globalmodel.TokenKey.
func (gs *globalService) GetUserIDFromContext(ctx context.Context) (int64, error) {
	userInfos, ok := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)
	if !ok || userInfos.ID == 0 {
		return 0, utils.BadRequest("invalid user context")
	}
	return userInfos.ID, nil
}
