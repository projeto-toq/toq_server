package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE id = ? AND deleted = 0;`

	entities, err := ua.Read(ctx, tx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_id.read_error", "error", err)
		// Propagate low-level error for service to translate
		return nil, err
	}

	if len(entities) == 0 {
		// Standard convention for not found at repository layer
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple users found with the same ID: %d", id)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_by_id.multiple_users_error", "user_id", id, "error", errMultiple)
		return nil, errMultiple
	}

	user, err = userconverters.UserEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_by_id.convert_error", "error", err)
		return nil, fmt.Errorf("convert user entity: %w", err)
	}

	// Note: Active role should be set by the calling service using Permission Service
	// This maintains separation of concerns between User and Permission domains

	return

}
