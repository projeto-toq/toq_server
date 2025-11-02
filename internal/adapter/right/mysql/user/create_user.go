package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) CreateUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	sql := `INSERT INTO users (
			full_name, nick_name, national_id, creci_number, creci_state, creci_validity,
			born_at, phone_number, email, zip_code, street, number, complement, neighborhood, 
			city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	entity := userconverters.UserDomainToEntity(user)

	result, execErr := ua.ExecContext(ctx, tx, "insert", sql,
		entity.FullName,
		entity.NickName,
		entity.NationalID,
		entity.CreciNumber,
		entity.CreciState,
		entity.CreciValidity,
		entity.BornAT,
		entity.PhoneNumber,
		entity.Email,
		entity.ZipCode,
		entity.Street,
		entity.Number,
		entity.Complement,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.Password,
		entity.OptStatus,
		entity.LastActivityAT,
		entity.Deleted,
		entity.LastSignInAttempt,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.create_user.exec_error", "error", execErr)
		return fmt.Errorf("create user: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.user.create_user.last_insert_id_error", "error", lastErr)
		return fmt.Errorf("user last insert id: %w", lastErr)
	}

	user.SetID(id)

	return
}
