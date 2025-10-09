package userservices

import (
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

func (us *userService) GetTokenExpiration(expired bool) int64 {
	if expired {
		return time.Now().UTC().Add(time.Hour * -1).Unix()
	}
	return time.Now().UTC().Add(globalmodel.GetAccessTTL()).Unix()
}
