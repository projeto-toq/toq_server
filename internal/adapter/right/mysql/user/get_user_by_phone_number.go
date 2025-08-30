package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUserByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entities, err := ua.Read(ctx, tx, "SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE phone_number = ?", phoneNumber)
	if err != nil {
		slog.Error("mysqluseradapter/GetUserByPhoneNumber: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter/GetUserByPhoneNumber: multiple users found with the same phone number", "phoneNumber", phoneNumber)
		return nil, utils.ErrInternalServer
	}

	user, err = userconverters.UserEntityToDomain(entities[0])
	if err != nil {
		return
	}

	// Note: Active role should be set by the calling service using Permission Service
	// This maintains separation of concerns between User and Permission domains

	return

}
