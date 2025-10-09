package mysqluseradapter

import (
	"database/sql"
	"errors"
	"fmt"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	devicetokenrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/device_token_repository"
)

// Ensure implementation satisfies port
var _ devicetokenrepository.DeviceTokenRepoPortInterface = (*DeviceTokenRepository)(nil)

type DeviceTokenRepository struct {
	db *sql.DB
}

func NewDeviceTokenRepository(db *sql.DB) *DeviceTokenRepository {
	return &DeviceTokenRepository{db: db}
}

func (r *DeviceTokenRepository) ListByUserID(userID int64) ([]usermodel.DeviceTokenInterface, error) {
	rows, err := r.db.Query(`SELECT id, user_id, device_token, device_id, platform FROM device_tokens WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []usermodel.DeviceTokenInterface
	for rows.Next() {
		var e userentity.DeviceTokenEntity
		var platform sql.NullString
		if err := rows.Scan(&e.ID, &e.UserID, &e.Token, &e.DeviceID, &platform); err != nil {
			return nil, err
		}
		dt := usermodel.NewDeviceToken()
		dt.SetID(e.ID)
		dt.SetUserID(e.UserID)
		dt.SetDeviceToken(e.Token)
		dt.SetDeviceID(e.DeviceID)
		result = append(result, dt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DeviceTokenRepository) AddToken(userID int64, token string, platform *string) (usermodel.DeviceTokenInterface, error) {
	if token == "" {
		return nil, errors.New("token required")
	}
	var platformArg any
	if platform != nil {
		platformArg = *platform
	}
	// Upsert-like: ignore duplicate token for same user
	res, err := r.db.Exec(`INSERT IGNORE INTO device_tokens (user_id, device_token, device_id, platform) VALUES (?,?,NULL,?)`, userID, token, platformArg)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	dt := usermodel.NewDeviceToken()
	dt.SetID(id)
	dt.SetUserID(userID)
	dt.SetDeviceToken(token)
	if platform != nil {
		// platform stored but not exposed in domain yet
	}
	return dt, nil
}

func (r *DeviceTokenRepository) RemoveToken(userID int64, token string) error {
	if token == "" {
		return errors.New("token required")
	}
	_, err := r.db.Exec(`DELETE FROM device_tokens WHERE user_id = ? AND device_token = ?`, userID, token)
	return err
}

func (r *DeviceTokenRepository) RemoveAllByUserID(userID int64) error {
	_, err := r.db.Exec(`DELETE FROM device_tokens WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("remove device tokens by user %d: %w", userID, err)
	}
	return nil
}

func (r *DeviceTokenRepository) ListTokensByOptedInUsers() ([]string, error) {
	rows, err := r.db.Query(`
		SELECT dt.device_token 
		FROM device_tokens dt 
		INNER JOIN users u ON dt.user_id = u.id 
		WHERE u.opt_status = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *DeviceTokenRepository) ListTokensByUserIDIfOptedIn(userID int64) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT dt.device_token 
		FROM device_tokens dt 
		INNER JOIN users u ON dt.user_id = u.id 
		WHERE dt.user_id = ? AND u.opt_status = 1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

// --- Per-device operations (schema-agnostic fallbacks) ---

// AddTokenForDevice falls back to AddToken when device_id column is not available.
func (r *DeviceTokenRepository) AddTokenForDevice(userID int64, deviceID, token string, platform *string) (usermodel.DeviceTokenInterface, error) {
	if token == "" {
		return nil, errors.New("token required")
	}
	if deviceID == "" {
		return nil, errors.New("deviceID required")
	}
	var platformArg any
	if platform != nil {
		platformArg = *platform
	}
	res, err := r.db.Exec(`INSERT INTO device_tokens (user_id, device_token, device_id, platform) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE device_token = VALUES(device_token), platform = VALUES(platform)`, userID, token, deviceID, platformArg)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	dt := usermodel.NewDeviceToken()
	dt.SetID(id)
	dt.SetUserID(userID)
	dt.SetDeviceToken(token)
	dt.SetDeviceID(deviceID)
	return dt, nil
}

// RemoveTokensByDeviceID removes tokens for a specific device; fallback removes none if device_id unsupported.
func (r *DeviceTokenRepository) RemoveTokensByDeviceID(userID int64, deviceID string) error {
	if deviceID == "" {
		return errors.New("deviceID required")
	}
	_, err := r.db.Exec(`DELETE FROM device_tokens WHERE user_id = ? AND device_id = ?`, userID, deviceID)
	return err
}

// ListTokensByDeviceID lists tokens for a specific device; fallback returns tokens by user if opted-in.
func (r *DeviceTokenRepository) ListTokensByDeviceID(userID int64, deviceID string) ([]string, error) {
	rows, err := r.db.Query(`SELECT device_token FROM device_tokens WHERE user_id = ? AND device_id = ?`, userID, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}
