package mysqluseradapter

import (
	"database/sql"
	"errors"
	"log/slog"

	userentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	devicetokenrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/device_token_repository"
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
	rows, err := r.db.Query(`SELECT id, user_id, device_token FROM device_tokens WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []usermodel.DeviceTokenInterface
	for rows.Next() {
		var e userentity.DeviceTokenEntity
		if err := rows.Scan(&e.ID, &e.UserID, &e.Token); err != nil {
			return nil, err
		}
		dt := usermodel.NewDeviceToken()
		dt.SetID(e.ID)
		dt.SetUserID(e.UserID)
		dt.SetDeviceToken(e.Token)
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
	// Upsert-like: ignore duplicate token for same user
	res, err := r.db.Exec(`INSERT IGNORE INTO device_tokens (user_id, device_token) VALUES (?,?)`, userID, token)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	dt := usermodel.NewDeviceToken()
	dt.SetID(id)
	dt.SetUserID(userID)
	dt.SetDeviceToken(token)
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
		slog.Warn("failed to remove device tokens", "userID", userID, "err", err)
	}
	return err
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
