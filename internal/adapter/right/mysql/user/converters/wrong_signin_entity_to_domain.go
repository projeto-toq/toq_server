package userconverters

import (
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func WrongSignInEntityToDomain(entity []any) (wsi usermodel.WrongSigninInterface, err error) {
	wsi = usermodel.NewWrongSignin()

	user_id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting user_id to int64", "value", entity[0])
		return nil, utils.ErrInternalServer
	}
	wsi.SetUserID(user_id)

	failed_attempts, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting failed_attempts to int64", "value", entity[1])
		return nil, utils.ErrInternalServer
	}
	wsi.SetFailedAttempts(failed_attempts)

	if entity[2] != nil {
		last_attempt_at, ok := entity[2].(time.Time)
		if !ok {
			slog.Error("Error converting last_attempt_at to time.Time", "value", entity[2])
			return nil, utils.ErrInternalServer
		}
		wsi.SetLastAttemptAt(last_attempt_at)
	}

	return
}
