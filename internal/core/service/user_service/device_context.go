package userservices

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// sanitizeDeviceContext normalizes and validates device token and device ID, updating context with the sanitized ID.
func (us *userService) sanitizeDeviceContext(ctx context.Context, deviceToken string, deviceID string, logPrefix string) (context.Context, string, string, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	trimmedToken := strings.TrimSpace(deviceToken)
	if trimmedToken == "" {
		err := utils.NewHTTPErrorWithSource(http.StatusBadRequest, "deviceToken is required")
		logger.Warn(logPrefix+".missing_device_token", "security", true)
		utils.SetSpanError(ctx, err)
		return ctx, "", "", err
	}

	trimmedDeviceID := strings.TrimSpace(deviceID)
	if trimmedDeviceID == "" {
		err := utils.NewHTTPErrorWithSource(http.StatusBadRequest, "deviceID is required")
		logger.Warn(logPrefix+".missing_device_id", "security", true)
		utils.SetSpanError(ctx, err)
		return ctx, "", "", err
	}

	if _, parseErr := uuid.Parse(trimmedDeviceID); parseErr != nil {
		err := utils.NewHTTPErrorWithSource(http.StatusBadRequest, "deviceID must be a valid UUID")
		logger.Warn(logPrefix+".invalid_device_id", "security", true, "device_id", trimmedDeviceID)
		utils.SetSpanError(ctx, err)
		return ctx, "", "", err
	}

	ctx = context.WithValue(ctx, globalmodel.DeviceIDKey, trimmedDeviceID)

	return ctx, trimmedToken, trimmedDeviceID, nil
}
