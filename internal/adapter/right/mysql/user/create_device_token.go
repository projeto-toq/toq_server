package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateDeviceToken inserts a new device token; ignores duplicate tokens for same user.
func (ua *UserAdapter) CreateDeviceToken(ctx context.Context, tx *sql.Tx, e userentity.DeviceTokenEntity) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Insert with device metadata when provided
	query := `INSERT INTO device_tokens (user_id, device_token, device_id, platform) VALUES (?, ?, ?, ?)`

	var deviceIDArg any
	if e.DeviceID != "" {
		deviceIDArg = e.DeviceID
	} else {
		deviceIDArg = nil
	}

	var platformArg any
	if e.Platform != nil {
		platformArg = *e.Platform
	} else {
		platformArg = nil
	}

	id, err := ua.Create(ctx, tx, query, e.UserID, e.Token, deviceIDArg, platformArg)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.create_device_token.create_error", "error", err)
		return 0, fmt.Errorf("insert device_token: %w", err)
	}
	return id, nil
}
