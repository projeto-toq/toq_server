package userservices

import (
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

func (us *userService) GetTokenExpiration(expired bool) int64 {
	if expired {
		return time.Now().UTC().Add(time.Hour * -1).Unix()
	}
	return time.Now().UTC().Add(globalmodel.TokenExpiration).Unix()
}
