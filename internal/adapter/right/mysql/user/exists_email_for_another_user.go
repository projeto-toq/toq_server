package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ExistsEmailForAnotherUser checks if an email is already used by a different user (deleted = 0)
func (ua *UserAdapter) ExistsEmailForAnotherUser(ctx context.Context, tx *sql.Tx, email string, excludeUserID int64) (bool, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT COUNT(id) as cnt FROM users WHERE email = ? AND id <> ? AND deleted = 0;`
	row := ua.QueryRowContext(ctx, tx, "select", query, email, excludeUserID)

	var cnt int64
	if scanErr := row.Scan(&cnt); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.exists_email_for_another.scan_error", "error", scanErr)
		return false, fmt.Errorf("exists email for another user scan: %w", scanErr)
	}

	return cnt > 0, nil
}
