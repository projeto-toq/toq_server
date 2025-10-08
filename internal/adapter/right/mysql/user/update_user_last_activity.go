package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserLastActivity(ctx context.Context, tx *sql.Tx, id int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE users SET last_activity_at = ? WHERE id = ?;`

	_, err = ua.Update(ctx, tx, query,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.update_user_last_activity.update_error", "error", err)
		return fmt.Errorf("update user last_activity: %w", err)
	}

	return
}
