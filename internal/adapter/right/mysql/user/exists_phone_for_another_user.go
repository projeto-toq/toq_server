package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
	entities, err := ua.Read(ctx, tx, query, phone, excludeUserID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.exists_phone_for_another.read_error", "error", err)
		return false, fmt.Errorf("exists phone for another user read: %w", err)
	}
	if len(entities) == 0 || len(entities[0]) == 0 {
		return false, nil
	}
	cnt, ok := entities[0][0].(int64)
	if !ok {
		errInvalid := fmt.Errorf("exists phone for another user: invalid count type %T", entities[0][0])
		utils.SetSpanError(ctx, errInvalid)
		logger.Error("mysql.user.exists_phone_for_another.invalid_count_type", "value", entities[0][0], "error", errInvalid)
		return false, errInvalid
	}
	return cnt > 0, nil
}
