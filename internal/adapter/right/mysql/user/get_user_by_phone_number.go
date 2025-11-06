package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUserByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE phone_number = ?`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, phoneNumber)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_by_phone.query_error", "error", queryErr)
		return nil, fmt.Errorf("get user by phone number query: %w", queryErr)
	}
	defer rows.Close()

	entities, err := scanUserEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_phone.scan_error", "error", err)
		return nil, fmt.Errorf("scan user by phone rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple users found for phone number: %s", phoneNumber)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_phone.multiple_users_error", "phone_number", phoneNumber, "error", errMultiple)
		return nil, errMultiple
	}

	user = userconverters.UserEntityToDomain(entities[0])

	return user, nil
}
