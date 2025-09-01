package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserPasswordByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `UPDATE users SET password = ? WHERE id = ?;`

	_, err = ua.Update(ctx, tx, query,
		user.GetPassword(),
		user.GetID(),
	)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserPasswordByID: error executing Update", "error", err)
		return fmt.Errorf("update user password: %w", err)
	}

	return
}
