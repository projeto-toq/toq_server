package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) VerifyUserDuplicity(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT count(id) as count
				FROM users WHERE (phone_number = ? OR email = ? OR national_id = ? ) AND deleted = 0;`

	entities, err := ua.Read(ctx, tx, query,
		user.GetPhoneNumber(),
		user.GetEmail(),
		user.GetNationalID(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.verify_duplicity.read_error", "error", err)
		return false, fmt.Errorf("verify user duplicity read: %w", err)
	}

	qty, ok := entities[0][0].(int64)
	if !ok {
		errInvalid := fmt.Errorf("verify user duplicity: invalid count type %T", entities[0][0])
		utils.SetSpanError(ctx, errInvalid)
		logger.Error("mysql.user.verify_duplicity.invalid_count_type", "value", entities[0][0], "error", errInvalid)
		return false, errInvalid
	}

	exist = qty > 0

	return

}
