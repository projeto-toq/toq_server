package mysqluseradapter

import (
	"context"
	"database/sql"

	mysqluseradapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
)

type UserAdapter struct {
	db *mysqluseradapter.Database
}

func NewUserAdapter(db *mysqluseradapter.Database) *UserAdapter {
	return &UserAdapter{
		db: db,
	}
}

// GetDeviceTokenRepository returns a repository bound to the underlying *sql.DB
// Avoids leaking *sql.DB outside adapter layer.
func (ua *UserAdapter) GetDeviceTokenRepository() *DeviceTokenRepository {
	return NewDeviceTokenRepository(ua.db.GetDB())
}

// AddDeviceToken adds a device token for a user (satisfies UserRepoPortInterface extension)
func (ua *UserAdapter) AddDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string, platform *string) error {
	if token == "" {
		return nil
	}
	e := userentity.DeviceTokenEntity{UserID: userID, Token: token}
	_, err := ua.CreateDeviceToken(ctx, tx, e)
	return err
}

// RemoveDeviceToken deletes a single device token for a user
func (ua *UserAdapter) RemoveDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string) error {
	if token == "" {
		return nil
	}
	_, err := tx.ExecContext(ctx, `DELETE FROM device_tokens WHERE user_id = ? AND device_token = ?`, userID, token)
	return err
}

func (ua *UserAdapter) RemoveAllDeviceTokens(ctx context.Context, tx *sql.Tx, userID int64) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM device_tokens WHERE user_id = ?`, userID)
	return err
}

// AddTokenForDevice stores a token associated to a device when schema supports it; fallback to user-only.
func (ua *UserAdapter) AddTokenForDevice(ctx context.Context, tx *sql.Tx, userID int64, deviceID, token string, platform *string) error {
	// Current schema lacks device_id column; fallback to AddDeviceToken
	return ua.AddDeviceToken(ctx, tx, userID, token, platform)
}

// RemoveTokensByDeviceID removes tokens for a specific device; no-op without device_id column.
func (ua *UserAdapter) RemoveTokensByDeviceID(ctx context.Context, tx *sql.Tx, userID int64, deviceID string) error {
	// No device_id in current schema; preserve data and do nothing
	return nil
}
