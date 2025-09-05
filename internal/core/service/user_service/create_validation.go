package userservices

import (
	"context"
	"database/sql"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (us *userService) CreateUserValidations(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {

	//set a fake validation as pending for email and phone to guarantee the profile is false. it will be replaced by the real validation afterwards
	validation := usermodel.NewValidation()
	validation.SetUserID(user.GetID())
	validation.SetNewEmail(user.GetEmail())
	validation.SetEmailCode(us.random6Digits())
	validation.SetEmailCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))
	validation.SetNewPhone(user.GetPhoneNumber())
	validation.SetPhoneCode(us.random6Digits())
	validation.SetPhoneCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))

	err = us.repo.UpdateUserValidations(ctx, tx, validation)
	if err != nil {
		// Método interno: infra  propagar com mapeamento alto-ndvel feito pelo chamador
		return err
	}

	return nil
}
