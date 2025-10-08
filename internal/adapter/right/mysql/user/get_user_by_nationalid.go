package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUserByNationalID(ctx context.Context, tx *sql.Tx, nationalID string) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entities, err := ua.Read(ctx, tx, "SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE national_id = ?", nationalID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_national_id.read_error", "error", err)
		return nil, fmt.Errorf("get user by national_id read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		errMultiple := errors.New("multiple users found for national_id")
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_national_id.multiple_users_error", "national_id", nationalID, "error", errMultiple)
		return nil, errMultiple
	}

	user, err = userconverters.UserEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_national_id.convert_error", "error", err)
		return nil, fmt.Errorf("convert user entity: %w", err)
	}

	// Note: Active role should be set by the calling service using Permission Service
	// This maintains separation of concerns between User and Permission domains

	return

}
