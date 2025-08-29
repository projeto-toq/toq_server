package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
"errors"
)

func (ua *UserAdapter) UpdateUserValidations(ctx context.Context, tx *sql.Tx, validation usermodel.ValidationInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Verify if the row should be blank at the end of the function and delete the row if it is
	if validation.GetEmailCode() == "" && validation.GetPhoneCode() == "" && validation.GetPasswordCode() == "" {
		_, err = ua.DeleteValidation(ctx, tx, validation.GetUserID())
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return
		}
	}

	sql := `INSERT INTO temp_user_validations (
		user_id, new_email, email_code, email_code_exp, 
		new_phone, phone_code, phone_code_exp, password_code, password_code_exp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		new_email = VALUES(new_email),
		email_code = VALUES(email_code),
		email_code_exp = VALUES(email_code_exp),
		new_phone = VALUES(new_phone),
		phone_code = VALUES(phone_code),
		phone_code_exp = VALUES(phone_code_exp),
		password_code = VALUES(password_code),
		password_code_exp = VALUES(password_code_exp);`

	entity := userconverters.UserValidationDomainToEntity(validation)

	id, err := ua.Update(ctx, tx, sql,
		entity.UserID,
		entity.NewEmail,
		entity.EmailCode,
		entity.EmailCodeExp,
		entity.NewPhone,
		entity.PhoneCode,
		entity.PhoneCodeExp,
		entity.PasswordCode,
		entity.PasswordCodeExp,
	)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserValidations: error executing Update", "error", err)
		return utils.ErrInternalServer
	}

	validation.SetUserID(id)

	return
}
