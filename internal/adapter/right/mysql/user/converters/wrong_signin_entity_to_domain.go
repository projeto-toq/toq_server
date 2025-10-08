package userconverters

import (
	"fmt"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func WrongSignInEntityToDomain(entity []any) (wsi usermodel.WrongSigninInterface, err error) {
	wsi = usermodel.NewWrongSignin()

	user_id, ok := entity[0].(int64)
	if !ok {
		return nil, fmt.Errorf("wrong signin converter: invalid user_id type %T", entity[0])
	}
	wsi.SetUserID(user_id)

	failed_attempts, ok := entity[1].(int64)
	if !ok {
		return nil, fmt.Errorf("wrong signin converter: invalid failed_attempts type %T", entity[1])
	}
	wsi.SetFailedAttempts(failed_attempts)

	if entity[2] != nil {
		last_attempt_at, ok := entity[2].(time.Time)
		if !ok {
			return nil, fmt.Errorf("wrong signin converter: invalid last_attempt_at type %T", entity[2])
		}
		wsi.SetLastAttemptAt(last_attempt_at)
	}

	return
}
