package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUsers(ctx context.Context, tx *sql.Tx) (users []usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE deleted = 0;`

	entities, err := ua.Read(ctx, tx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_users.read_error", "error", err)
		return nil, fmt.Errorf("get users read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, entity := range entities {
		user, convertErr := userconverters.UserEntityToDomain(entity)
		if convertErr != nil {
			utils.SetSpanError(ctx, convertErr)
			logger.Error("mysql.user.get_users.convert_error", "error", convertErr)
			return nil, fmt.Errorf("convert user entity: %w", convertErr)
		}

		// Note: Active role should be set by the calling service using Permission Service
		// This maintains separation of concerns between User and Permission domains

		users = append(users, user)
	}

	return

}
