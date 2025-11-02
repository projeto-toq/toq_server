package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserValidations(ctx context.Context, tx *sql.Tx, validation usermodel.ValidationInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// If all validation codes are empty, cleanup the row instead of upserting blanks
	if validation.GetEmailCode() == "" && validation.GetPhoneCode() == "" && validation.GetPasswordCode() == "" {
		_, err = ua.DeleteValidation(ctx, tx, validation.GetUserID())
		return err
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

	result, execErr := ua.ExecContext(ctx, tx, "insert", sql,
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
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_validations.exec_error", "error", execErr)
		return fmt.Errorf("update user validations: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_validations.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user validations update rows affected: %w", rowsErr)
	}

	validation.SetUserID(entity.UserID)

	return
}
