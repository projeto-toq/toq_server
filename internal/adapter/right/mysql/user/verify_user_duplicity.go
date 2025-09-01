package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) VerifyUserDuplicity(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT count(id) as count
				FROM users WHERE (phone_number = ? OR email = ? OR national_id = ? ) AND deleted = 0;`

	entities, err := ua.Read(ctx, tx, query,
		user.GetPhoneNumber(),
		user.GetEmail(),
		user.GetNationalID(),
	)
	if err != nil {
		slog.Error("mysqluseradapter/VerifyUserDuplicity: error executing Read", "error", err)
		return false, fmt.Errorf("verify user duplicity read: %w", err)
	}

	qty, ok := entities[0][0].(int64)
	if !ok {
		slog.Error("mysqluseradapter/VerifyUserDuplicity: error converting qty to int64", "qty", entities[0][0])
		return false, fmt.Errorf("verify user duplicity: invalid count type %T", entities[0][0])
	}

	exist = qty > 0

	return

}
