package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// AddDeviceToken registers a new push notification token for a user device using UPSERT semantics
//
// This function uses INSERT ... ON DUPLICATE KEY UPDATE to handle token rotation:
//   - If device_id doesn't exist: creates new record
//   - If device_id exists: updates token and platform (handles token refresh)
//
// The database schema requires a UNIQUE constraint on (user_id, device_id) for this to work correctly.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging propagation
//   - tx: Database transaction (can be nil for standalone operation)
//   - userID: User's unique identifier (foreign key to users.id)
//   - deviceID: Unique device identifier (UUIDv4 format)
//   - token: FCM or APNs push notification token (up to 255 characters)
//   - platform: Device platform pointer ("android"/"ios"/"web") - optional (NULL if nil)
//
// Returns:
//   - token: Created or updated device token record with ID populated
//   - error: Database errors (connection, constraint violations, etc.)
//
// Business Rules:
//   - One token per device per user (enforced by UNIQUE constraint)
//   - Token rotation supported (ON DUPLICATE KEY UPDATE)
//   - Platform is optional and can be updated
//
// Example:
//
//	token, err := adapter.AddDeviceToken(ctx, tx, 123, "550e8400-e29b-41d4-a716-446655440000", "fcm_abc...", &platform)
func (ua *UserAdapter) AddDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, deviceID, token string, platform *string) (usermodel.DeviceToken, error) {
	// Initialize tracing for observability (metrics + distributed tracing)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return usermodel.DeviceToken{}, err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// UPSERT query: insert new record or update existing token on duplicate device_id
	// Requires UNIQUE constraint on (user_id, device_id) in schema
	query := `INSERT INTO device_tokens (user_id, device_id, device_token, platform) 
	          VALUES (?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE 
	              device_token = VALUES(device_token), 
	              platform = VALUES(platform)`

	// Handle NULL platform (convert Go nil to SQL NULL)
	var platformArg interface{}
	if platform != nil {
		platformArg = *platform
	}

	// Execute INSERT/UPDATE using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "insert", query, userID, deviceID, token, platformArg)
	if execErr != nil {
		// Mark span as error for distributed tracing analysis
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.add_device_token.exec_error", "user_id", userID, "device_id", deviceID, "error", execErr)
		return usermodel.DeviceToken{}, fmt.Errorf("add device token: %w", execErr)
	}

	// Get auto-incremented ID (note: ON DUPLICATE KEY UPDATE returns existing ID)
	id, idErr := result.LastInsertId()
	if idErr != nil {
		utils.SetSpanError(ctx, idErr)
		logger.Error("mysql.user.add_device_token.last_insert_id_error", "user_id", userID, "device_id", deviceID, "error", idErr)
		return usermodel.DeviceToken{}, fmt.Errorf("get last insert id: %w", idErr)
	}

	// Build entity to convert to domain Value Object
	entity := userentity.DeviceTokenEntity{
		ID:       id,
		UserID:   userID,
		Token:    token,
		DeviceID: deviceID,
		Platform: sql.NullString{
			String: func() string {
				if platform != nil {
					return *platform
				}
				return ""
			}(),
			Valid: platform != nil,
		},
	}

	// Convert database entity to domain Value Object
	return userconverters.DeviceTokenEntityToVO(entity), nil
}
