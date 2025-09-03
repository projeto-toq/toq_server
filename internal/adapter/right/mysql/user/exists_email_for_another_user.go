package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ExistsEmailForAnotherUser checks if an email is already used by a different user (deleted = 0)
func (ua *UserAdapter) ExistsEmailForAnotherUser(ctx context.Context, tx *sql.Tx, email string, excludeUserID int64) (bool, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	query := `SELECT COUNT(id) as cnt FROM users WHERE email = ? AND id <> ? AND deleted = 0;`
	entities, err := ua.Read(ctx, tx, query, email, excludeUserID)
	if err != nil {
		slog.Error("mysqluseradapter/ExistsEmailForAnotherUser: error executing Read", "error", err)
		return false, fmt.Errorf("exists email for another user read: %w", err)
	}
	if len(entities) == 0 || len(entities[0]) == 0 {
		return false, nil
	}
	cnt, ok := entities[0][0].(int64)
	if !ok {
		slog.Error("mysqluseradapter/ExistsEmailForAnotherUser: invalid count type", "value", entities[0][0])
		return false, fmt.Errorf("exists email for another user: invalid count type %T", entities[0][0])
	}
	return cnt > 0, nil
}
