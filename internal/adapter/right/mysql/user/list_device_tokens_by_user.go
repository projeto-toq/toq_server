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

// ListDeviceTokensByUserID retrieves all push notification tokens for a user
//
// This function returns complete device token records (with ID, platform, device_id).
// Returns empty slice if user has no tokens (NOT sql.ErrNoRows).
//
// Parameters:
//   - ctx: Context for tracing and logging propagation
//   - tx: Database transaction (can be nil for standalone operation)
//   - userID: User's unique identifier
//
// Returns:
//   - tokens: Slice of DeviceToken Value Objects (empty if no tokens)
//   - error: Database errors (connection, query execution, scan failures)
//
// Business Rules:
//   - Returns all tokens regardless of user's opt_status or deleted status
//   - Returns complete token records (not just token strings)
//   - Empty result is NOT an error (user may have no registered devices)
//
// Use Cases:
//   - Admin operations: view all user devices
//   - Device management: list user's registered devices
//   - Debugging: inspect token state
//
// Example:
//
//	tokens, err := adapter.ListDeviceTokensByUserID(ctx, nil, 123)
//	// len(tokens) == 0 if user has no devices (not an error)
func (ua *UserAdapter) ListDeviceTokensByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]usermodel.DeviceToken, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query lists all columns explicitly (NEVER use SELECT *)
	query := `SELECT id, user_id, device_token, device_id, platform 
	          FROM device_tokens 
	          WHERE user_id = ?`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, userID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.list_device_tokens.query_error", "user_id", userID, "error", queryErr)
		return nil, fmt.Errorf("list device tokens: %w", queryErr)
	}
	defer rows.Close()

	// Scan rows into entities
	var entities []userentity.DeviceTokenEntity
	for rows.Next() {
		var entity userentity.DeviceTokenEntity
		scanErr := rows.Scan(&entity.ID, &entity.UserID, &entity.Token, &entity.DeviceID, &entity.Platform)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.user.list_device_tokens.scan_error", "user_id", userID, "error", scanErr)
			return nil, fmt.Errorf("scan device token: %w", scanErr)
		}
		entities = append(entities, entity)
	}

	// Check for iteration errors
	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.list_device_tokens.rows_error", "user_id", userID, "error", rowsErr)
		return nil, fmt.Errorf("iterate device tokens: %w", rowsErr)
	}

	// Convert entities to domain Value Objects
	return userconverters.DeviceTokenEntitiesToVOs(entities), nil
}
