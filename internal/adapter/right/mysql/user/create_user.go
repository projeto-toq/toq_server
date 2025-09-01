package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) CreateUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO users (
			full_name, nick_name, national_id, creci_number, creci_state, creci_validity,
			born_at, phone_number, email, zip_code, street, number, complement, neighborhood, 
			city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	entity := userconverters.UserDomainToEntity(user)

	id, err := ua.Create(ctx, tx, sql,
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
	if err != nil {
		slog.Error("mysqluseradapter/CreateUser: error executing Create", "error", err)
		return fmt.Errorf("create user: %w", err)
	}

	user.SetID(id)

	return
}
