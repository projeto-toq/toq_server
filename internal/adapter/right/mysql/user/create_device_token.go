package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	userentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/entities"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateDeviceToken inserts a new device token; ignores duplicate tokens for same user.
func (ua *UserAdapter) CreateDeviceToken(ctx context.Context, tx *sql.Tx, e userentity.DeviceTokenEntity) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	// Table currently has columns: id, user_id, device_token
	query := `INSERT INTO device_tokens (user_id, device_token) VALUES (?, ?)`
	id, err := ua.Create(ctx, tx, query, e.UserID, e.Token)
	if err != nil {
		slog.Error("mysqluseradapter/CreateDeviceToken: insert failed", "err", err)
		return 0, fmt.Errorf("insert device_token: %w", err)
	}
	return id, nil
}
