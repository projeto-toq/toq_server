package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserPasswordByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE users SET password = ? WHERE id = ?;`

	_, err = ua.Update(ctx, tx, query,
		user.GetPassword(),
		user.GetID(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.update_user_password.update_error", "error", err)
		return fmt.Errorf("update user password: %w", err)
	}

	return
}
