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

	row := ua.QueryRowContext(ctx, tx, "select", query,
		user.GetPhoneNumber(),
		user.GetEmail(),
		user.GetNationalID(),
	)

	var qty int64
	if scanErr := row.Scan(&qty); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.verify_duplicity.scan_error", "error", scanErr)
		return false, fmt.Errorf("verify user duplicity scan: %w", scanErr)
	}

	exist = qty > 0

	return

}
