package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ExistsPhoneForAnotherUser checks if a phone is already used by a different user (deleted = 0)
func (ua *UserAdapter) ExistsPhoneForAnotherUser(ctx context.Context, tx *sql.Tx, phone string, excludeUserID int64) (bool, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT COUNT(id) as cnt FROM users WHERE phone_number = ? AND id <> ? AND deleted = 0;`
	row := ua.QueryRowContext(ctx, tx, "select", query, phone, excludeUserID)

	var cnt int64
	if scanErr := row.Scan(&cnt); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.user.exists_phone_for_another.scan_error", "error", scanErr)
		return false, fmt.Errorf("exists phone for another user scan: %w", scanErr)
	}

	return cnt > 0, nil
}
