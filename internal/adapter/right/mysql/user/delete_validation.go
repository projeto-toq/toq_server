package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteValidation(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM temp_user_validations WHERE user_id = ?;`

	deleted, err = ua.Delete(ctx, tx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.delete_validation.delete_error", "error", err)
		return 0, fmt.Errorf("delete validation: %w", err)
	}

	// Idempotent: if no rows were deleted, that's fine (nothing to clean up)
	return deleted, nil
}
